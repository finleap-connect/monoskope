package certificatemanagement

import (
	"github.com/google/uuid"
)

// CertificateManager is an interface to read/write certificate requests and certificates
type CertificateManager interface {
	// RequestCertificate creates a new CertificateRequest with the given
	// base64 encoded string of a PEM encoded certificate request.
	RequestCertificate(requestID uuid.UUID, csr []byte) error
}
