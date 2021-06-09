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
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/k8s"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
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
		expectedIssuer := "monoskope-issuer"
		expectedDuration := time.Hour * 48
		expectedCSR := []byte("some-csr-bytes")

		testEnv, err := jwt.NewTestEnv(test.NewTestEnv("TestReactors"))
		Expect(err).NotTo(HaveOccurred())
		defer util.PanicOnError(testEnv.Shutdown())

		When("ClusterCreated event occurs", func() {
			aggregateId := uuid.New()
			aggregateVersion := uint64(1)
			eventType := events.ClusterCreated
			eventData := &eventdata.ClusterCreated{
				Name:                "TestCluster",
				Label:               "test-cluster",
				ApiServerAddress:    "https://localhost",
				CaCertificateBundle: []byte("somecabundle"),
			}

			It("emits a ClusterBootstrapTokenCreated event", func() {
				eventChannel := make(chan eventsourcing.Event, 1)

				k8sClient := k8s.NewMockClient(mockCtrl)
				reactor := NewClusterBootstrapReactor(testEnv.CreateSigner(), certificatemanagement.NewCertManagerClient(k8sClient, expectedNamespace, expectedIssuer, expectedDuration))

				err := reactor.HandleEvent(ctx, eventsourcing.NewEvent(ctx, eventType, eventsourcing.ToEventDataFromProto(eventData), time.Now().UTC(), aggregateType, aggregateId, aggregateVersion), eventChannel)
				Expect(err).NotTo(HaveOccurred())

				event := <-eventChannel
				Expect(event.EventType()).To(Equal(events.ClusterBootstrapTokenCreated))

				eventDataTokenCreated := &eventdata.ClusterBootstrapTokenCreated{}
				err = event.Data().ToProto(eventDataTokenCreated)
				Expect(err).NotTo(HaveOccurred())
				Expect(eventDataTokenCreated.JWT).To(Not(BeEmpty()))
			})
		})
		When("ClusterCertificateRequested event occurs", func() {
			aggregateId := uuid.New()
			aggregateVersion := uint64(2)
			eventType := events.ClusterCertificateRequested
			eventData := &eventdata.ClusterCertificateRequested{
				CertificateSigningRequest: expectedCSR,
			}

			cr := new(cmapi.CertificateRequest)
			cr.Spec.Usages = append(cr.Spec.Usages, cmapi.UsageClientAuth)
			cr.Spec.IssuerRef.Kind = cmapi.IssuerKind
			cr.Spec.IssuerRef.Group = cmapi.IssuerGroupAnnotationKey
			cr.Spec.IsCA = false
			cr.Name = aggregateId.String()
			cr.Namespace = expectedNamespace
			cr.Spec.Request = expectedCSR
			cr.Spec.IssuerRef.Name = expectedIssuer
			cr.Spec.Duration = &v1.Duration{
				Duration: expectedDuration,
			}

			It("emits a ClusterOperatorCertificateRequestIssued event", func() {
				eventChannel := make(chan eventsourcing.Event, 2)

				k8sClient := k8s.NewMockClient(mockCtrl)
				reactor := NewClusterBootstrapReactor(testEnv.CreateSigner(), certificatemanagement.NewCertManagerClient(k8sClient, expectedNamespace, expectedIssuer, expectedDuration))
				expectedCACert := []byte("some-ca-cert")
				expectedCert := []byte("some-cert")

				k8sClient.EXPECT().Get(ctx, types.NamespacedName{Name: aggregateId.String(), Namespace: expectedNamespace}, gomock.Any()).
					Return(errors.NewNotFound(cmapi.Resource(cr.Name), cr.Name))

				k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, obj runtime.Object) error {
					cr := obj.(*cmapi.CertificateRequest)
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
					Expect(event.EventType()).To(Equal(events.ClusterOperatorCertificateRequestIssued))

					event = <-eventChannel
					Expect(event).NotTo(BeNil())
					Expect(event.EventType()).To(Equal(events.ClusterOperatorCertificateIssued))
				}()

				err := reactor.HandleEvent(ctx, eventsourcing.NewEvent(ctx, eventType, eventsourcing.ToEventDataFromProto(eventData), time.Now().UTC(), aggregateType, aggregateId, aggregateVersion), eventChannel)
				Expect(err).NotTo(HaveOccurred())

				<-finished
			})
		})
	})
})
