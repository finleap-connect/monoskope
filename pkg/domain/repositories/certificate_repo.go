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

package repositories

import (
	"context"

	domApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	esErrors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
)

type certificateRepository struct {
	es.Repository
}

// CertificateRepository is a repository for reading and writing certificate projections.
type CertificateRepository interface {
	es.Repository
	ReadOnlyCertificateRepository
	WriteOnlyCertificateRepository
}

// ReadOnlyCertificateRepository is a repository for reading certificate projections.
type ReadOnlyCertificateRepository interface {
	// GetCertificate retrieves certificates by aggregate type and id
	GetCertificate(context.Context, *domApi.GetCertificateRequest) (*projections.Certificate, error)
}

// WriteOnlyCertificateRepository is a repository for writing certificate projections.
type WriteOnlyCertificateRepository interface {
}

// NewCertificateRepository creates a repository for reading and writing certificate projections.
func NewCertificateRepository(repository es.Repository) CertificateRepository {
	return &certificateRepository{
		Repository: repository,
	}
}

// Retrieve certificates for a specified aggregate ID and type.
func (r *certificateRepository) GetCertificate(ctx context.Context, req *domApi.GetCertificateRequest) (*projections.Certificate, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	for _, projection := range ps {
		certificate, ok := projection.(*projections.Certificate)
		if !ok {
			return nil, esErrors.ErrInvalidProjectionType
		}
		if certificate.ReferencedAggregateId == req.AggregateId && certificate.AggregateType == req.AggregateType {
			return certificate, nil
		}
	}

	return nil, esErrors.ErrProjectionNotFound
}
