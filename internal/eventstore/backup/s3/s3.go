package s3

import (
	"archive/tar"
	"context"
	"crypto/aes"
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

func (b *S3BackupHandler) RunBackup(ctx context.Context) (*backup.BackupResult, error) {
	filename := fmt.Sprintf("%v-eventstore-backup.tar", time.Now().UTC().Format(time.RFC3339))
	b.log.Info("Starting backup...", "Bucket", b.conf.Bucket, "Endpoint", b.conf.Bucket, "Filename", filename)

	result := &backup.BackupResult{}
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

func (b *S3BackupHandler) streamEvents(ctx context.Context, writer *io.PipeWriter, result *backup.BackupResult) error {
	defer writer.Close()
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	eventStream, err := b.store.Load(ctx, &es.StoreQuery{})
	if err != nil {
		return err
	}

	encryptionKey := ""
	if v := os.Getenv("S3_ENCRYPTION_KEY"); v != "" {
		encryptionKey = v
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

		bytes, err := json.Marshal(event)
		if err != nil {
			b.log.Error(err, "An error occured when marshalling event", "AggregateType", event.AggregateType())
			return err
		}

		// Use encryption if key has been specified
		if encryptionKey != "" {
			encryptedBytes, err := encryptAES([]byte(encryptionKey), bytes)
			if err != nil {
				return err
			}
			bytes = encryptedBytes
		}

		err = tarWriter.WriteHeader(&tar.Header{
			Name:       fmt.Sprintf("%s-%s-%v.json", event.AggregateType().String(), event.AggregateID().String(), event.AggregateVersion()),
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

func encryptAES(key []byte, payload []byte) ([]byte, error) {
	// create cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// allocate space for ciphered data
	out := make([]byte, len(payload))

	// encrypt
	c.Encrypt(out, payload)

	// return encrypted payload
	return out, nil
}

// func decryptAES(key []byte, encryptedPayload []byte) ([]byte, error) {
// 	c, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	payload := make([]byte, len(encryptedPayload))
// 	c.Decrypt(payload, encryptedPayload)

// 	return payload, nil
// }
