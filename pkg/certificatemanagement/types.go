package certificatemanagement

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrRequestFailed  = errors.New("requesting certificate failed")
	ErrRequestInvalid = errors.New("certificate request invalid")
	ErrRequestPending = errors.New("certificate request pending")
	ErrRequestDenied  = errors.New("certificate request denied")
)

// CertificateManager is an interface to read/write certificate requests and certificates
type CertificateManager interface {
	// RequestCertificate creates a new CertificateRequest with the given
	// base64 encoded string of a PEM encoded certificate request.
	RequestCertificate(ctx context.Context, requestID uuid.UUID, csr []byte) error
	// GetCertificate returns a byte slice containing a PEM encoded signed certificate
	// resulting from a certificate signing request identified by the requestID
	GetCertificate(ctx context.Context, requestID uuid.UUID) ([]byte, error)
}
