// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s3

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/finleap-connect/monoskope/internal/test"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/storage"
	"github.com/ory/dockertest/v3"
)

const (
	accessKeyID     = "TESTACCESSKEY"
	secretAccessKey = "TESTSECRETKEY"
	encryptionKey   = "thisis32bitlongpassphraseimusing"
)

type TestEnv struct {
	*test.TestEnv
	storageTestEnv *storage.TestEnv
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

	err = env.storageTestEnv.Store.Open(context.Background())
	if err != nil {
		return nil, err
	}

	if err := env.CreateDockerPool(false); err != nil {
		return nil, err
	}

	container, err := env.Run(&dockertest.RunOptions{
		Name:       "minio",
		Repository: "minio/minio",
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

	err = env.CreateBucket(env.Endpoint, accessKeyID, secretAccessKey)
	if err != nil {
		return nil, err
	}

	os.Setenv("S3_ENCRYPTION_KEY", encryptionKey)
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
}
