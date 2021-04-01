package backup

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ory/dockertest/v3"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/storage"
)

const (
	accessKeyID     = "TESTACCESSKEY"
	secretAccessKey = "TESTSECRETKEY"
)

type TestEnv struct {
	*test.TestEnv
	storageTestEnv *storage.TestEnv
	store          es.Store
	Endpoint       string
	Bucket         string
}

func NewTestEnvWithParent(testEnv *test.TestEnv) (*TestEnv, error) {
	var err error

	env := &TestEnv{
		TestEnv: testEnv,
		Bucket:  "test-bucket",
	}

	env.storageTestEnv, err = storage.NewTestEnvWithParent(testEnv)
	if err != nil {
		return nil, err
	}

	env.store, err = storage.NewPostgresEventStore(env.storageTestEnv.GetStoreConfig())
	if err != nil {
		return nil, err
	}
	err = env.store.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	if err := env.CreateDockerPool(); err != nil {
		return nil, err
	}

	if v := os.Getenv("MINIO_URL"); v != "" {
		env.Endpoint = v
	} else {
		container, err := env.Run(&dockertest.RunOptions{
			Name:       "minio",
			Repository: "artifactory.figo.systems/public_docker/minio/minio",
			Tag:        "latest",
			Cmd:        []string{"server", "/data"},
			Env: []string{
				fmt.Sprintf("MINIO_ACCESS_KEY=%s", accessKeyID),
				fmt.Sprintf("MINIO_SECRET_KEY=%s", secretAccessKey),
			},
		})
		if err != nil {
			return nil, err
		}

		env.Endpoint = fmt.Sprintf("localhost:%s", container.GetPort("9000/tcp"))

		env.Log.Info("check minio connection", "endpoint", env.Endpoint)
		err = env.WaitForS3(env.Endpoint, accessKeyID, secretAccessKey)
		if err != nil {
			return nil, err
		}
	}

	err = env.CreateBucket(env.Endpoint, accessKeyID, secretAccessKey)
	if err != nil {
		return nil, err
	}

	os.Setenv("S3_ACCESS_KEY", accessKeyID)
	os.Setenv("S3_SECRET_KEY", secretAccessKey)

	return env, nil
}

func (t *TestEnv) WaitForS3(endpoint, accessKeyID, secretAccessKey string) error {
	return t.Retry(func() error {
		newSession, err := session.NewSession(&aws.Config{
			Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
			Endpoint:         aws.String(endpoint),
			Region:           aws.String("us-east-1"),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		client := s3.New(newSession)
		input := &s3.ListBucketsInput{}
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		_, err = client.ListBucketsWithContext(ctx, input)
		return err
	})
}

func (t *TestEnv) CreateBucket(endpoint, accessKeyID, secretAccessKey string) error {
	return t.Retry(func() error {
		newSession, err := session.NewSession(&aws.Config{
			Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
			Endpoint:         aws.String(endpoint),
			Region:           aws.String("us-east-1"),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		client := s3.New(newSession)
		_, err = client.CreateBucket(&s3.CreateBucketInput{Bucket: &t.Bucket})
		return err
	})
}
