package backup

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var _ = Describe("s3 backups", func() {
	It("should backup s3 bucket", func() {
		conf := &S3Config{
			Bucket:     testEnv.Bucket,
			Endpoint:   testEnv.Endpoint,
			DisableSSL: true,
			Region:     "us-east-1",
		}

		var err error
		userId := uuid.New()
		roleBindingId := uuid.New()
		err = testEnv.store.Save(context.Background(), []eventsourcing.Event{
			eventsourcing.NewEvent(context.Background(), events.UserCreated, nil, time.Now().UTC(), aggregates.User, userId, 1),
		})
		Expect(err).ToNot(HaveOccurred())
		err = testEnv.store.Save(context.Background(), []eventsourcing.Event{
			eventsourcing.NewEvent(context.Background(), events.UserRoleBindingCreated, nil, time.Now().UTC(), aggregates.UserRoleBinding, roleBindingId, 1),
			eventsourcing.NewEvent(context.Background(), events.UserRoleBindingDeleted, nil, time.Now().UTC(), aggregates.UserRoleBinding, roleBindingId, 2),
		})
		Expect(err).ToNot(HaveOccurred())

		b := NewS3BackupHandler(conf, testEnv.store)
		Expect(b).ToNot(BeNil())

		backupResult, err := b.RunBackup(context.Background())
		Expect(err).ToNot(HaveOccurred())
		Expect(backupResult.ProcessedEvents).To(BeNumerically(">", 0))
		Expect(backupResult.ProcessedBytes).To(BeNumerically(">", 0))
	})
})
