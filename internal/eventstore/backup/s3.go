package backup

import (
	"context"
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
)

const encryptionAlgorithm string = "AES256"
const partitionSize int = 1000

type S3Config struct {
	endpoint       string
	bucket         string
	retentionCount int
	disableSSL     bool
}

type S3BackupHandler struct {
	log           logger.Logger
	store         es.Store
	conf          *S3Config
	s3Client      *s3.S3
	encryptionKey string
}

func NewBackupHandler(conf *S3Config, store es.Store) (*S3BackupHandler, error) {
	b := &S3BackupHandler{
		log:   logger.WithName("s3-backup-handler"),
		conf:  conf,
		store: store,
	}
	return b, nil
}

func (b *S3BackupHandler) initClient() error {
	if b.s3Client != nil {
		return nil
	}

	dstAccessKey := os.Getenv("S3_ACCESS_KEY_DST")
	dstSecretKey := os.Getenv("S3_SECRET_KEY_DST")

	if v := os.Getenv("S3_ENCRYPTION_KEY_DST"); v != "" {
		b.encryptionKey = v
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(dstAccessKey, dstSecretKey, ""),
		Endpoint:         &b.conf.endpoint,
		Region:           aws.String("us-east-1"),
		DisableSSL:       &b.conf.disableSSL,
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return err
	}
	b.s3Client = s3.New(sess)

	return nil
}

func (b *S3BackupHandler) RunBackup(ctx context.Context) (int64, error) {
	filename := fmt.Sprintf("%v-eventstore-backup.tar", time.Now().UTC().Format(time.RFC3339))

	b.log.Info("Starting backup...", "Bucket", b.conf.bucket, "Endpoint", b.conf.bucket, "Filename", filename)
	err := b.initClient()
	if err != nil {
		return 0, err
	}
	return b.backupEvents(ctx, filename)
}

func (b *S3BackupHandler) backupEvents(ctx context.Context, filename string) (int64, error) {
	return 0, nil
}

func (b *S3BackupHandler) uploadBackup(ctx context.Context, filename string, reader *io.PipeReader, size int64) error {
	b.log.Info("Uploading backup to S3...", "Bucket", b.conf.bucket, "Endpoint", b.conf.bucket, "Filename", filename)

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploaderWithClient(b.s3Client)

	// Prevent hitting the limit
	if size/uploader.PartSize > int64(uploader.MaxUploadParts) {
		uploader.PartSize = int64(float64(size) / float64(uploader.MaxUploadParts) * 1.2) // * <factor> for safety if initial object size increases during backup
		b.log.Info("Adjusted s3manager.Uploader.PartSize to prevent hitting s3manager.Uploader.MaxUploadParts", "partSize", uploader.PartSize, "MaxUploadParts", uploader.MaxUploadParts)
	}

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
	if err != nil {
		return err
	}
	return nil
}
