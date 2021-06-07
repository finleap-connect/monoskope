package certificatemanagement

import (
	"time"

	"github.com/google/uuid"
	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	"k8s.io/client-go/kubernetes"
)

type certManagerClient struct {
	issuer        string
	namingPattern string
	namespace     string
	duration      time.Duration
}

// NewCertManagerClient creates a cert-manager.io specific implementation of the certificatemanagement.CertificateManager interface
func NewCertManagerClient(k8sClient *kubernetes.Clientset, namingPattern, namespace, issuer string, duration time.Duration) CertificateManager {
	return &certManagerClient{
		issuer:        issuer,
		namingPattern: namingPattern,
		namespace:     namespace,
		duration:      duration,
	}
}

// RequestCertificate requests a new certificate based on the given certificate signing request
func (c *certManagerClient) RequestCertificate(requestID uuid.UUID, csr []byte) error {
	cr := &cmapi.CertificateRequest{}

	// fixed defaults
	cr.Spec.Usages = append(cr.Spec.Usages, "server auth", "client auth")
	cr.Spec.IssuerRef.Kind = cmapi.IssuerKind
	cr.Spec.IssuerRef.Group = cmapi.IssuerGroupAnnotationKey
	cr.Spec.IsCA = false

	// input
	cr.Name = requestID.String()
	cr.Namespace = c.namespace
	cr.Spec.Request = csr
	cr.Spec.IssuerRef.Name = c.issuer
	cr.Spec.Duration.Duration = c.duration

	return nil
}
