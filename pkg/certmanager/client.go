package certmanager

import (
	"k8s.io/client-go/kubernetes"
)

// CertManagerClient is an interface to read/write cert-manager resources
type CertManagerClient interface {
	CreateCertificateRequest(namespace, name string, csr []byte) error
}

type certManagerClient struct {
}

func NewClient(k8sClient *kubernetes.Clientset) CertManagerClient {
	return &certManagerClient{}
}

func (c *certManagerClient) CreateCertificateRequest(namespace, name string, csr []byte) error {
	return nil
}
