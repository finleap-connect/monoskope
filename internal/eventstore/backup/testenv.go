package backup

import (
	"context"
	"fmt"
	"os"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
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
	dstEndpoint    string
}

func NewTestEnv() (*TestEnv, error) {
	var err error
	env := &TestEnv{
		TestEnv: test.NewTestEnv("EventStoreTestEnv"),
	}

	env.storageTestEnv, err = storage.NewTestEnv()
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

	// Start rabbitmq
	container, err := env.Run(&dockertest.RunOptions{
		Name:       "minio",
		Repository: "artifactory.figo.systems/public_docker/minio/minio",
		Tag:        "RELEASE.2020-10-03T02-54-56Z",
		Cmd:        []string{"server", "/data"},
		PortBindings: map[dc.Port][]dc.PortBinding{
			"9000": {{HostPort: "9000"}},
		},
		Env: []string{
			fmt.Sprintf("MINIO_ACCESS_KEY=%s", accessKeyID),
			fmt.Sprintf("MINIO_SECRET_KEY=%s", secretAccessKey),
		},
	})
	if err != nil {
		return nil, err
	}

	env.dstEndpoint = fmt.Sprintf("localhost:%s", container.GetPort("9000/tcp"))
	env.Log.Info("check minio connection", "endpoint", env.dstEndpoint)

	os.Setenv("S3_ACCESS_KEY", accessKeyID)
	os.Setenv("S3_SECRET_KEY", secretAccessKey)

	return env, nil
}
