package s3

import (
	"archive/tar"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/backup"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type S3Config struct {
	Endpoint   string
	Bucket     string
	Region     string
	DisableSSL bool
}

type S3BackupHandler struct {
	log      logger.Logger
	store    es.Store
	conf     *S3Config
	s3Client *s3.S3
}

func NewS3BackupHandler(conf *S3Config, store es.Store) backup.BackupHandler {
	b := &S3BackupHandler{
		log:   logger.WithName("s3-backup-handler"),
		conf:  conf,
		store: store,
	}
	return b
}

func (b *S3BackupHandler) initClient() error {
	if b.s3Client != nil {
		return nil
	}

	dstAccessKey := os.Getenv("S3_ACCESS_KEY")
	dstSecretKey := os.Getenv("S3_SECRET_KEY")

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(dstAccessKey, dstSecretKey, ""),
		Endpoint:         aws.String(b.conf.Endpoint),
		Region:           aws.String(b.conf.Region),
		DisableSSL:       aws.Bool(b.conf.DisableSSL),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return err
	}
	b.s3Client = s3.New(sess)

	return nil
}

func (b *S3BackupHandler) RunBackup(ctx context.Context) (*backup.Result, error) {
	filename := fmt.Sprintf("%v-eventstore-backup.tar", time.Now().UTC().Format(time.RFC3339))
	b.log.Info("Starting backup...", "Bucket", b.conf.Bucket, "Endpoint", b.conf.Bucket, "Filename", filename)

	result := &backup.Result{BackupIdentifier: filename}
	err := b.initClient()
	if err != nil {
		return result, err
	}

	reader, writer := io.Pipe()
	var eg errgroup.Group
	eg.Go(func() error {
		return b.streamEvents(ctx, writer, result)
	})
	eg.Go(func() error {
		return b.uploadBackup(ctx, filename, reader)
	})

	if err := eg.Wait(); err != nil {
		b.log.Error(err, "Error occured when backing up eventstore.", "ProcessedEvents", result.ProcessedEvents, "ProcessedBytes", result.ProcessedBytes)
		return result, err
	}
	b.log.Info("Backing up eventstore has been successful.", "ProcessedEvents", result.ProcessedEvents, "ProcessedBytes", result.ProcessedBytes)
	return result, nil
}

func (b *S3BackupHandler) RunRestore(ctx context.Context, identifier string) (*backup.Result, error) {
	b.log.Info("Starting restore...", "Bucket", b.conf.Bucket, "Endpoint", b.conf.Bucket, "Filename", identifier)

	result := &backup.Result{BackupIdentifier: identifier}
	err := b.initClient()
	if err != nil {
		return result, err
	}

	reader, writer := io.Pipe()
	var eg errgroup.Group
	eg.Go(func() error {
		return b.readBackup(ctx, writer, identifier)
	})
	eg.Go(func() error {
		return b.storeEvents(ctx, reader, result)
	})

	if err := eg.Wait(); err != nil {
		b.log.Error(err, "Error occured when restoring events.", "ProcessedEvents", result.ProcessedEvents, "ProcessedBytes", result.ProcessedBytes)
		return result, err
	}
	b.log.Info("Restoring events has been successful.", "ProcessedEvents", result.ProcessedEvents, "ProcessedBytes", result.ProcessedBytes)
	return result, nil
}

func (b *S3BackupHandler) storeEvents(ctx context.Context, reader *io.PipeReader, result *backup.Result) error {
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	var encryptionKey []byte
	if v := os.Getenv("S3_ENCRYPTION_KEY"); v != "" {
		encryptionKey = []byte(v)
		b.log.Info("Decrypting backup with key specified in env var S3_ENCRYPTION_KEY.")
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		bytes := make([]byte, header.Size)
		n, err := io.ReadFull(tarReader, bytes)
		if err != nil {
			return err
		}

		// Use encryption if key has been specified
		if len(encryptionKey) > 0 {
			encryptedBytes, err := decryptAES(encryptionKey, bytes)
			if err != nil {
				return err
			}
			bytes = encryptedBytes
		}

		event := &s3Event{}
		err = json.Unmarshal(bytes, event)
		if err != nil {
			b.log.Error(err, "An error occured when unmarshalling event", "AggregateType", event.AggregateType())
			return err
		}

		err = b.store.Save(ctx, []es.Event{event})
		if err != nil {
			return err
		}

		result.ProcessedBytes += uint64(n)
		result.ProcessedEvents++
	}
	return nil
}

func (b *S3BackupHandler) streamEvents(ctx context.Context, writer *io.PipeWriter, result *backup.Result) error {
	defer writer.Close()
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	eventStream, err := b.store.Load(ctx, &es.StoreQuery{})
	if err != nil {
		return err
	}

	var encryptionKey []byte
	if v := os.Getenv("S3_ENCRYPTION_KEY"); v != "" {
		encryptionKey = []byte(v)
		b.log.Info("Encrypting backup with AES and key specified in env var S3_ENCRYPTION_KEY.")
	}

	b.log.Info("Streaming events from store...")
	for {
		event, err := eventStream.Receive()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		bytes, err := json.Marshal(convertToS3Event(event))
		if err != nil {
			b.log.Error(err, "An error occured when marshalling event", "AggregateType", event.AggregateType())
			return err
		}

		// Use encryption if key has been specified
		if len(encryptionKey) > 0 {
			encryptedBytes, err := encryptAES(encryptionKey, bytes)
			if err != nil {
				return err
			}
			bytes = encryptedBytes
		}

		err = tarWriter.WriteHeader(&tar.Header{
			Name:       fmt.Sprintf("%v", result.ProcessedEvents),
			Mode:       0600,
			ChangeTime: time.Now().UTC(),
			ModTime:    time.Now().UTC(),
			Size:       int64(len(bytes)),
		})
		if err != nil {
			b.log.Error(err, "An error occured when writing tar header", "AggregateType", event.AggregateType())
			return err
		}

		numBytes, err := tarWriter.Write(bytes)
		if err != nil {
			b.log.Error(err, "An error occured when writing tar payload for event", "AggregateType", event.AggregateType())
			return err
		} else {
			result.ProcessedEvents++
			result.ProcessedBytes += uint64(numBytes)
		}
	}

	return nil
}

func (b *S3BackupHandler) uploadBackup(ctx context.Context, filename string, reader *io.PipeReader) error {
	defer reader.Close()
	b.log.Info("Uploading backup to S3...", "Bucket", b.conf.Bucket, "Endpoint", b.conf.Bucket, "Filename", filename)

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploaderWithClient(b.s3Client)

	ui := &s3manager.UploadInput{
		Bucket:      aws.String(b.conf.Bucket),
		Key:         aws.String(filename),
		Body:        reader,
		ContentType: aws.String("application/tar"),
	}

	_, err := uploader.UploadWithContext(ctx, ui)
	return err
}

func (b *S3BackupHandler) readBackup(ctx context.Context, writer *io.PipeWriter, filename string) error {
	objectInput := &s3.GetObjectInput{
		Bucket: aws.String(b.conf.Bucket),
		Key:    aws.String(filename),
	}

	object, err := b.s3Client.GetObjectWithContext(ctx, objectInput)
	if err != nil {
		b.log.Error(err, "An error occured when backing up object")
		return err
	}

	_, err = io.Copy(writer, object.Body)
	if err != nil {
		b.log.Error(err, "An error occured when writing object to tar")
		return err
	}

	return nil
}

func encryptAES(key []byte, plaintext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	// return encrypted payload
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptAES(key []byte, ciphertext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func convertToS3Event(event es.Event) es.Event {
	return &s3Event{
		EType:      event.EventType(),
		EData:      event.Data(),
		ETimestamp: event.Timestamp(),
		AType:      event.AggregateType(),
		AID:        event.AggregateID(),
		AVersion:   event.AggregateVersion(),
		MD:         event.Metadata(),
	}
}

// s3Event is the private implementation of the Event interface for a s3 event backup event.
type s3Event struct {
	EType      es.EventType
	EData      es.EventData
	ETimestamp time.Time
	AType      es.AggregateType
	AID        uuid.UUID
	AVersion   uint64
	MD         map[string]string
}

// EventType implements the EventType method of the Event interface.
func (e s3Event) EventType() es.EventType {
	return e.EType
}

// Data implements the Data method of the Event interface.
func (e s3Event) Data() es.EventData {
	return e.EData
}

// Timestamp implements the Timestamp method of the Event interface.
func (e s3Event) Timestamp() time.Time {
	return e.ETimestamp
}

// AggregateType implements the AggregateType method of the Event interface.
func (e s3Event) AggregateType() es.AggregateType {
	return e.AType
}

// AggrgateID implements the AggrgateID method of the Event interface.
func (e s3Event) AggregateID() uuid.UUID {
	return e.AID
}

// AggregateVersion implements the AggregateVersion method of the Event interface.
func (e s3Event) AggregateVersion() uint64 {
	return e.AVersion
}

// AggregateVersion implements the AggregateVersion method of the Event interface.
func (e s3Event) Metadata() map[string]string {
	return e.MD
}

// String implements the String method of the Event interface.
func (e s3Event) String() string {
	return fmt.Sprintf("%s@%d", e.EventType(), e.AggregateVersion())
}
