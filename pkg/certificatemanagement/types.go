// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	GetCertificate(ctx context.Context, requestID uuid.UUID) (ca []byte, cert []byte, err error)
}
