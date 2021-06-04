package certificatemanagement

import (
	"time"

	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	"k8s.io/client-go/kubernetes"
)

type certManagerClient struct {
}

func NewCertManagerClient(k8sClient *kubernetes.Clientset) CertificateManager {
	return &certManagerClient{}
}

func (c *certManagerClient) CreateCertificateRequest(name, namespace, issuer string, request []byte, duration time.Duration) error {
	cr := &cmapi.CertificateRequest{}

	// fixed defaults
	cr.Spec.Usages = append(cr.Spec.Usages, "server auth", "client auth")
	cr.Spec.IssuerRef.Kind = cmapi.IssuerKind
	cr.Spec.IssuerRef.Group = cmapi.IssuerGroupAnnotationKey
	cr.Spec.IsCA = false

	// input
	cr.Name = name
	cr.Namespace = namespace
	cr.Spec.Request = request
	cr.Spec.IssuerRef.Name = issuer
	cr.Spec.Duration.Duration = duration

	return nil
}
