package certificatemanagement

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/k8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("package certificatemanagement", func() {
	Context("CertManagerClient", func() {
		var mockCtrl *gomock.Controller
		ctx := context.Background()
		expectedCSRID := uuid.New()
		expectedNamespace := "monoskope"
		expectedIssuer := "monoskope-issuer"
		expectedDuration := time.Hour * 48

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		When("RequestCertificate is called with a valid CSR", func() {
			expectedCSR := []byte("some-csr-bytes")

			It("returns no error", func() {
				k8sClient := k8s.NewMockClient(mockCtrl)

				cr := new(cmapi.CertificateRequest)
				cr.Spec.Usages = append(cr.Spec.Usages, "server auth", "client auth")
				cr.Spec.IssuerRef.Kind = cmapi.IssuerKind
				cr.Spec.IssuerRef.Group = cmapi.IssuerGroupAnnotationKey
				cr.Spec.IsCA = false
				cr.Name = expectedCSRID.String()
				cr.Namespace = expectedNamespace
				cr.Spec.Request = expectedCSR
				cr.Spec.IssuerRef.Name = expectedIssuer
				cr.Spec.Duration = &v1.Duration{
					Duration: expectedDuration,
				}

				k8sClient.EXPECT().Create(ctx, cr).Return(nil)

				client := NewCertManagerClient(k8sClient, expectedNamespace, expectedIssuer, expectedDuration)

				err := client.RequestCertificate(ctx, expectedCSRID, expectedCSR)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("GetCertificate is called", func() {
			expectedCert := []byte("some-cert")

			It("returns the issued cert with no error", func() {
				k8sClient := k8s.NewMockClient(mockCtrl)
				k8sClient.EXPECT().Get(ctx, types.NamespacedName{Name: expectedCSRID.String(), Namespace: expectedNamespace}, new(cmapi.CertificateRequest)).DoAndReturn(func(_ context.Context, _ types.NamespacedName, obj runtime.Object) error {
					cr := obj.(*cmapi.CertificateRequest)
					cr.Status.Conditions = append(cr.Status.Conditions, cmapi.CertificateRequestCondition{
						Type:    cmapi.CertificateRequestConditionReady,
						Message: "Certificate ready.",
					})
					cr.Status.Certificate = expectedCert
					k8sClient.EXPECT().Delete(ctx, cr).Return(nil)
					return nil
				})

				client := NewCertManagerClient(k8sClient, expectedNamespace, expectedIssuer, expectedDuration)
				cert, err := client.GetCertificate(ctx, expectedCSRID)
				Expect(err).NotTo(HaveOccurred())
				Expect(cert).To(Equal(expectedCert))
			})

			// Checks the GetCertificate method returns the right errors based on the condition the CertificateRequest is in
			checkErrorResponse := func(condition cmapi.CertificateRequestConditionType, expectedError error) {
				k8sClient := k8s.NewMockClient(mockCtrl)
				k8sClient.EXPECT().Get(ctx, types.NamespacedName{Name: expectedCSRID.String(), Namespace: expectedNamespace}, new(cmapi.CertificateRequest)).DoAndReturn(func(_ context.Context, _ types.NamespacedName, obj runtime.Object) error {
					cr := obj.(*cmapi.CertificateRequest)
					cr.Status.Conditions = append(cr.Status.Conditions, cmapi.CertificateRequestCondition{
						Type:    condition,
						Message: string(condition),
					})
					return nil
				})

				client := NewCertManagerClient(k8sClient, expectedNamespace, expectedIssuer, expectedDuration)
				cert, err := client.GetCertificate(ctx, expectedCSRID)
				Expect(expectedError).To(HaveOccurred())
				Expect(cert).To(BeNil())
				Expect(err).To(Equal(expectedError))
			}

			It("returns ErrRequestPending when condition is CertificateRequestConditionApproved", func() {
				checkErrorResponse(cmapi.CertificateRequestConditionApproved, ErrRequestPending)
			})
			It("returns ErrRequestInvalid when condition is CertificateRequestConditionInvalidRequest", func() {
				checkErrorResponse(cmapi.CertificateRequestConditionInvalidRequest, ErrRequestInvalid)
			})
			It("returns ErrRequestDenied when condition is CertificateRequestConditionDenied", func() {
				checkErrorResponse(cmapi.CertificateRequestConditionDenied, ErrRequestDenied)
			})
		})
	})
})
