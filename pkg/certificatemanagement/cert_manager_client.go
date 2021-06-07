package certificatemanagement

import (
	"context"
	"time"

	"github.com/google/uuid"
	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/k8s"
)

type certManagerClient struct {
	k8sClient     k8s.K8sClient
	issuer        string
	namingPattern string
	namespace     string
	duration      time.Duration
}

// NewCertManagerClient creates a cert-manager.io specific implementation of the certificatemanagement.CertificateManager interface
func NewCertManagerClient(k8sClient k8s.K8sClient, namingPattern, namespace, issuer string, duration time.Duration) CertificateManager {
	return &certManagerClient{
		k8sClient:     k8sClient,
		issuer:        issuer,
		namingPattern: namingPattern,
		namespace:     namespace,
		duration:      duration,
	}
}

// RequestCertificate requests a new certificate based on the given certificate signing request
func (c *certManagerClient) RequestCertificate(ctx context.Context, requestID uuid.UUID, csr []byte) error {
	cr := new(cmapi.CertificateRequest)

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

	return c.k8sClient.Create(ctx, cr)
}
