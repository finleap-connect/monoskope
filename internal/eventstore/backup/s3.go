package backup

import (
	"archive/tar"
	"context"
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
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"golang.org/x/sync/errgroup"
)

const encryptionAlgorithm string = "AES256"

type S3Config struct {
	endpoint   string
	bucket     string
	region     string
	disableSSL bool
}

type S3BackupHandler struct {
	log           logger.Logger
	store         es.Store
	conf          *S3Config
	s3Client      *s3.S3
	encryptionKey string
}

func NewS3BackupHandler(conf *S3Config, store es.Store) BackupHandler {
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

	if v := os.Getenv("S3_ENCRYPTION_KEY"); v != "" {
		b.encryptionKey = v
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(dstAccessKey, dstSecretKey, ""),
		Endpoint:         aws.String(b.conf.endpoint),
		Region:           aws.String(b.conf.region),
		DisableSSL:       aws.Bool(b.conf.disableSSL),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return err
	}
	b.s3Client = s3.New(sess)

	return nil
}

func (b *S3BackupHandler) RunBackup(ctx context.Context) (*BackupResult, error) {
	filename := fmt.Sprintf("%v-eventstore-backup.tar", time.Now().UTC().Format(time.RFC3339))
	b.log.Info("Starting backup...", "Bucket", b.conf.bucket, "Endpoint", b.conf.bucket, "Filename", filename)

	result := &BackupResult{}
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
	return result, eg.Wait()
}

func (b *S3BackupHandler) streamEvents(ctx context.Context, writer *io.PipeWriter, result *BackupResult) error {
	defer writer.Close()
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	eventStream, err := b.store.Load(ctx, &es.StoreQuery{})
	if err != nil {
		return err
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
	b.log.Info("Uploading backup to S3...", "Bucket", b.conf.bucket, "Endpoint", b.conf.bucket, "Filename", filename)

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploaderWithClient(b.s3Client)

	ui := &s3manager.UploadInput{
		Bucket:      aws.String(b.conf.bucket),
		Key:         aws.String(filename),
		Body:        reader,
		ContentType: aws.String("application/tar"),
	}
	if b.encryptionKey != "" {
		ui.SSECustomerKey = &b.encryptionKey
		ui.SSECustomerAlgorithm = aws.String(encryptionAlgorithm)
	}

	_, err := uploader.UploadWithContext(ctx, ui)
	return err
}
