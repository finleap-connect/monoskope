package certificatemanagement

import (
	"context"
	"time"

	"github.com/google/uuid"
	apiutil "github.com/jetstack/cert-manager/pkg/api/util"
	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type certManagerClient struct {
	log       logger.Logger
	k8sClient ctrlclient.Client
	issuer    string
	namespace string
	duration  time.Duration
}

// NewCertManagerClient creates a cert-manager.io specific implementation of the certificatemanagement.CertificateManager interface
func NewCertManagerClient(k8sClient ctrlclient.Client, namespace, issuer string, duration time.Duration) CertificateManager {
	return &certManagerClient{
		log:       logger.WithName("certManagerClient"),
		k8sClient: k8sClient,
		issuer:    issuer,
		namespace: namespace,
		duration:  duration,
	}
}

// RequestCertificate requests a new certificate based on the given certificate signing request
func (c *certManagerClient) RequestCertificate(ctx context.Context, requestID uuid.UUID, csr []byte) error {
	cr := new(cmapi.CertificateRequest)

	// fixed defaults
	cr.Spec.Usages = append(cr.Spec.Usages, cmapi.UsageClientAuth)
	cr.Spec.IssuerRef.Kind = cmapi.IssuerKind
	cr.Spec.IssuerRef.Group = cmapi.IssuerGroupAnnotationKey
	cr.Spec.IsCA = false

	// input
	cr.Name = requestID.String()
	cr.Namespace = c.namespace
	cr.Spec.Request = csr
	cr.Spec.IssuerRef.Name = c.issuer
	cr.Spec.Duration = &v1.Duration{
		Duration: c.duration,
	}

	c.log.Info("Requesting certificate...", "RequestID", requestID.String(), "Namespace", c.namespace, "Issuer", c.issuer)
	err := c.k8sClient.Create(ctx, cr)
	if err != nil {
		c.log.Error(err, "Requesting certificate failed.", "RequestID", requestID.String(), "Namespace", c.namespace, "Issuer", c.issuer)
		return ErrRequestFailed
	}
	return nil
}

// GetCertificate returns a byte slice containing a PEM encoded signed certificate resulting from a certificate signing request identified by the requestID
func (c *certManagerClient) GetCertificate(ctx context.Context, requestID uuid.UUID) ([]byte, []byte, error) {
	cr := new(cmapi.CertificateRequest)
	err := c.k8sClient.Get(ctx, types.NamespacedName{Name: requestID.String(), Namespace: c.namespace}, cr)
	if err != nil {
		return nil, nil, err
	}

	if apiutil.CertificateRequestHasInvalidRequest(cr) {
		return nil, nil, ErrRequestInvalid
	}
	if apiutil.CertificateRequestIsDenied(cr) {
		return nil, nil, ErrRequestDenied
	}

	if len(cr.Status.Certificate) > 0 {
		if err := c.k8sClient.Delete(ctx, cr); err != nil {
			c.log.Error(err, "Failed to delete request after successfull certificate issueing.", "RequestID", requestID.String(), "Namespace", c.namespace, "Issuer", c.issuer)
		}
		return cr.Status.CA, cr.Status.Certificate, nil
	}

	return nil, nil, ErrRequestPending
}
