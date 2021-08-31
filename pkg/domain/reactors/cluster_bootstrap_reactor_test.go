// Copyright 2021 Monoskope Authors
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

package reactors

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	apiutil "github.com/jetstack/cert-manager/pkg/api/util"
	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/certificatemanagement"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	mock_k8s "gitlab.figo.systems/platform/monoskope/monoskope/test/k8s"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("package reactors", func() {
	var (
		mockCtrl *gomock.Controller
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("ClusterBootstrapReactor", func() {
		ctx := context.Background()
		aggregateType := aggregates.Cluster
		expectedNamespace := "monoskope"
		expectedIssuer := "selfsigning-issuer"
		expectedIssuerKind := "ClusterIssuer"
		expectedDuration := time.Hour * 48
		expectedCSR := []byte("some-csr-bytes")
		expectedIssuerUrl := "https://localhost"

		testEnv, err := jwt.NewTestEnv(test.NewTestEnv("TestReactors"))
		Expect(err).NotTo(HaveOccurred())
		defer util.PanicOnError(testEnv.Shutdown())

		When("ClusterCreated event occurs", func() {
			aggregateId := uuid.New()
			aggregateVersion := uint64(1)
			eventType := events.ClusterCreatedV2
			eventData := &eventdata.ClusterCreatedV2{
				DisplayName:         "TestCluster",
				Name:                "test-cluster",
				ApiServerAddress:    "https://localhost",
				CaCertificateBundle: []byte("somecabundle"),
			}

			It("emits a ClusterBootstrapTokenCreated,UserCreated,UserRoleBindingCreated event", func() {
				eventChannel := make(chan eventsourcing.Event, 3)

				k8sClient := mock_k8s.NewMockClient(mockCtrl)
				reactor := NewClusterBootstrapReactor(expectedIssuerUrl, testEnv.CreateSigner(), certificatemanagement.NewCertManagerClient(k8sClient, expectedNamespace, expectedIssuerKind, expectedIssuer, expectedDuration))

				err := reactor.HandleEvent(ctx, eventsourcing.NewEvent(ctx, eventType, eventsourcing.ToEventDataFromProto(eventData), time.Now().UTC(), aggregateType, aggregateId, aggregateVersion), eventChannel)
				Expect(err).NotTo(HaveOccurred())

				event := <-eventChannel
				Expect(event.EventType()).To(Equal(events.ClusterBootstrapTokenCreated))

				eventDataTokenCreated := &eventdata.ClusterBootstrapTokenCreated{}
				err = event.Data().ToProto(eventDataTokenCreated)
				Expect(err).NotTo(HaveOccurred())
				Expect(eventDataTokenCreated.Jwt).To(Not(BeEmpty()))

				ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
				defer cancel()

				select {
				case <-ctxWithTimeout.Done():
					Fail("Context deadline exceeded waiting for event.")
				case event = <-eventChannel:
					Expect(event.EventType()).To(Equal(events.UserCreated))
				}

				select {
				case <-ctxWithTimeout.Done():
					Fail("Context deadline exceeded waiting for event.")
				case event = <-eventChannel:
					Expect(event.EventType()).To(Equal(events.UserRoleBindingCreated))
				}
			})
		})
		When("CertificateRequested event occurs", func() {
			aggregateId := uuid.New()
			aggregateVersion := uint64(2)
			eventType := events.CertificateRequested
			eventData := &eventdata.CertificateRequested{
				SigningRequest:          expectedCSR,
				ReferencedAggregateId:   aggregateId.String(),
				ReferencedAggregateType: aggregates.Cluster.String(),
			}

			cr := new(cmapi.CertificateRequest)
			cr.Spec.Usages = append(cr.Spec.Usages, cmapi.UsageClientAuth)
			cr.Spec.IsCA = false
			cr.Name = aggregateId.String()
			cr.Namespace = expectedNamespace
			cr.Spec.Request = expectedCSR
			cr.Spec.IssuerRef.Kind = expectedIssuerKind
			cr.Spec.IssuerRef.Name = expectedIssuer
			cr.Spec.Duration = &v1.Duration{
				Duration: expectedDuration,
			}

			It("emits a CertificateRequestIssued event", func() {
				eventChannel := make(chan eventsourcing.Event, 2)

				k8sClient := mock_k8s.NewMockClient(mockCtrl)
				reactor := NewClusterBootstrapReactor(expectedIssuerUrl, testEnv.CreateSigner(), certificatemanagement.NewCertManagerClient(k8sClient, expectedNamespace, expectedIssuerKind, expectedIssuer, expectedDuration))
				expectedCACert := []byte("some-ca-cert")
				expectedCert := []byte("some-cert")

				k8sClient.EXPECT().Get(gomock.Any(), types.NamespacedName{Name: aggregateId.String(), Namespace: expectedNamespace}, gomock.Any()).
					Return(errors.NewNotFound(cmapi.Resource(cr.Name), cr.Name))

				k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, obj runtime.Object) error {
					cr := obj.(*cmapi.CertificateRequest)
					k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ types.NamespacedName, obj runtime.Object) error {
						crGet := obj.(*cmapi.CertificateRequest)
						*crGet = *cr
						apiutil.SetCertificateRequestCondition(crGet, cmapi.CertificateRequestConditionApproved, cmmeta.ConditionTrue, "Approved by test.", "Certificate approved.")
						return nil
					})
					k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ types.NamespacedName, obj runtime.Object) error {
						crGet := obj.(*cmapi.CertificateRequest)
						*crGet = *cr
						apiutil.SetCertificateRequestCondition(crGet, cmapi.CertificateRequestConditionReady, cmmeta.ConditionTrue, "Approved by test.", "Certificate ready.")
						crGet.Status.Certificate = expectedCert
						crGet.Status.CA = expectedCACert
						return nil
					})
					return nil
				})
				k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

				finished := make(chan bool)
				go func() {
					defer GinkgoRecover()
					defer func() { finished <- true }()

					event := <-eventChannel
					Expect(event).NotTo(BeNil())
					Expect(event.EventType()).To(Equal(events.CertificateRequestIssued))

					event = <-eventChannel
					Expect(event).NotTo(BeNil())
					Expect(event.EventType()).To(Equal(events.CertificateIssued))
				}()

				err := reactor.HandleEvent(ctx, eventsourcing.NewEvent(ctx, eventType, eventsourcing.ToEventDataFromProto(eventData), time.Now().UTC(), aggregateType, aggregateId, aggregateVersion), eventChannel)
				Expect(err).NotTo(HaveOccurred())

				<-finished
			})
		})
	})
})
