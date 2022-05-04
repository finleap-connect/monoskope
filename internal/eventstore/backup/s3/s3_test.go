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
	"time"

	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("s3", func() {
	backupIdentifier := ""

	conf := &S3Config{
		DisableSSL: true,
		Region:     "us-east-1",
	}

	It("should backup eventstore", func() {
		defer testEnv.storageTestEnv.ClearStore(context.Background())
		conf.Bucket = testEnv.Bucket
		conf.Endpoint = testEnv.Endpoint

		var err error
		userId := uuid.New()
		roleBindingId := uuid.New()

		err = testEnv.storageTestEnv.Store.Save(context.Background(), []eventsourcing.Event{
			eventsourcing.NewEvent(context.Background(), events.UserCreated, nil, time.Now().UTC(), aggregates.User, userId, 1),
		})
		Expect(err).ToNot(HaveOccurred())

		err = testEnv.storageTestEnv.Store.Save(context.Background(), []eventsourcing.Event{
			eventsourcing.NewEvent(context.Background(), events.UserRoleBindingCreated, nil, time.Now().UTC(), aggregates.UserRoleBinding, roleBindingId, 1),
			eventsourcing.NewEvent(context.Background(), events.UserRoleBindingDeleted, nil, time.Now().UTC(), aggregates.UserRoleBinding, roleBindingId, 2),
		})
		Expect(err).ToNot(HaveOccurred())

		b := NewS3BackupHandler(conf, testEnv.storageTestEnv.Store, 0)
		Expect(b).ToNot(BeNil())

		result, err := b.RunBackup(context.Background())
		Expect(err).ToNot(HaveOccurred())
		Expect(result.ProcessedEvents).To(BeNumerically(">", 0))
		Expect(result.ProcessedBytes).To(BeNumerically(">", 0))
		backupIdentifier = result.BackupIdentifier
		time.Sleep(1000 * time.Millisecond)
	})
	It("should restore eventstore", func() {
		defer testEnv.storageTestEnv.ClearStore(context.Background())
		conf.Bucket = testEnv.Bucket
		conf.Endpoint = testEnv.Endpoint

		b := NewS3BackupHandler(conf, testEnv.storageTestEnv.Store, 0)
		Expect(b).ToNot(BeNil())

		result, err := b.RunRestore(context.Background(), backupIdentifier)
		Expect(err).ToNot(HaveOccurred())
		Expect(result.ProcessedEvents).To(BeNumerically(">", 0))
		Expect(result.ProcessedBytes).To(BeNumerically(">", 0))
	})
	It("should purge backups", func() {
		defer testEnv.storageTestEnv.ClearStore(context.Background())
		conf.Bucket = testEnv.Bucket
		conf.Endpoint = testEnv.Endpoint

		var err error
		userId := uuid.New()
		roleBindingId := uuid.New()

		err = testEnv.storageTestEnv.Store.Save(context.Background(), []eventsourcing.Event{
			eventsourcing.NewEvent(context.Background(), events.UserCreated, nil, time.Now().UTC(), aggregates.User, userId, 1),
		})
		Expect(err).ToNot(HaveOccurred())

		err = testEnv.storageTestEnv.Store.Save(context.Background(), []eventsourcing.Event{
			eventsourcing.NewEvent(context.Background(), events.UserRoleBindingCreated, nil, time.Now().UTC(), aggregates.UserRoleBinding, roleBindingId, 1),
			eventsourcing.NewEvent(context.Background(), events.UserRoleBindingDeleted, nil, time.Now().UTC(), aggregates.UserRoleBinding, roleBindingId, 2),
		})
		Expect(err).ToNot(HaveOccurred())

		b := NewS3BackupHandler(conf, testEnv.storageTestEnv.Store, 5)
		Expect(b).ToNot(BeNil())

		for i := 0; i < 8; i++ {
			result, err := b.RunBackup(context.Background())
			Expect(err).ToNot(HaveOccurred())
			Expect(result.ProcessedEvents).To(BeNumerically(">", 0))
			Expect(result.ProcessedBytes).To(BeNumerically(">", 0))
		}

		pr, err := b.RunPurge(context.Background())
		Expect(err).ToNot(HaveOccurred())
		Expect(pr.BackupsLeft).To(BeNumerically("==", 5))
		Expect(pr.PurgedBackups).To(BeNumerically("==", 4))
	})
})
