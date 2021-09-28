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
	"time"

	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/google/uuid"
	apiutil "github.com/jetstack/cert-manager/pkg/api/util"
	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type certManagerClient struct {
	log        logger.Logger
	k8sClient  ctrlclient.Client
	issuer     string
	issuerKind string
	namespace  string
	duration   time.Duration
}

// NewCertManagerClient creates a cert-manager.io specific implementation of the certificatemanagement.CertificateManager interface
func NewCertManagerClient(k8sClient ctrlclient.Client, namespace, issuerKind, issuer string, duration time.Duration) CertificateManager {
	return &certManagerClient{
		log:        logger.WithName("certManagerClient"),
		k8sClient:  k8sClient,
		issuer:     issuer,
		issuerKind: issuerKind,
		namespace:  namespace,
		duration:   duration,
	}
}

// RequestCertificate requests a new certificate based on the given certificate signing request
func (c *certManagerClient) RequestCertificate(ctx context.Context, requestID uuid.UUID, csr []byte) error {
	cr := new(cmapi.CertificateRequest)

	if err := c.k8sClient.Get(ctx, types.NamespacedName{Name: requestID.String(), Namespace: c.namespace}, cr); err != nil {
		if apierrors.IsNotFound(err) {
			// fixed defaults
			cr.Spec.Usages = append(cr.Spec.Usages, cmapi.UsageClientAuth)
			cr.Spec.IsCA = false

			// input
			cr.Name = requestID.String()
			cr.Namespace = c.namespace
			cr.Spec.Request = csr
			cr.Spec.IssuerRef.Name = c.issuer
			cr.Spec.IssuerRef.Kind = c.issuerKind
			cr.Spec.Duration = &v1.Duration{
				Duration: c.duration,
			}
		} else {
			return err
		}
	}

	c.log.Info("Requesting certificate...", "RequestID", requestID.String(), "Namespace", c.namespace, "Issuer", c.issuer)
	if err := c.k8sClient.Create(ctx, cr); err != nil {
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
