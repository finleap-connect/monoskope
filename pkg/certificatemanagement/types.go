package certificatemanagement

import "time"

// CertificateManager is an interface to read/write certificate requests and certificates
type CertificateManager interface {
	// CreateCertificateRequest creates a new CertificateRequest with the given
	// base64 encoded string of a PEM encoded certificate request.
	CreateCertificateRequest(name, namespace, issuer string, request []byte, duration time.Duration) error
}
