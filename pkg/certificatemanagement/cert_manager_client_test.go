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
)

var _ = Describe("package certificatemanagement", func() {
	Context("CertManagerClient", func() {
		var mockCtrl *gomock.Controller
		ctx := context.Background()

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		When("RequestCertificate is called with a valid CSR", func() {
			expectedCSRID := uuid.New()
			expectedCSR := []byte("some-csr-bytes")
			expectedNamespace := "monoskope"
			expectedIssuer := "monoskope-issuer"
			expectedDuration := time.Hour * 48

			It("no error occurs", func() {
				k8sClient := k8s.NewMockK8sClient(mockCtrl)

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
	})
})
