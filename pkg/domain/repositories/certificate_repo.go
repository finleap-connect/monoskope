// Copyright 2022 Monoskope Authors
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
	DomainRepository[*projections.Certificate]
}

// CertificateRepository is a repository for reading and writing certificate projections.
type CertificateRepository interface {
	DomainRepository[*projections.Certificate]
	// GetCertificate retrieves certificates by aggregate type and id
	GetCertificate(context.Context, *domApi.GetCertificateRequest) (*projections.Certificate, error)
}

// NewCertificateRepository creates a repository for reading and writing certificate projections.
func NewCertificateRepository(repository es.Repository[*projections.Certificate]) CertificateRepository {
	return &certificateRepository{
		NewDomainRepository(repository),
	}
}

// Retrieve certificates for a specified aggregate ID and type.
func (r *certificateRepository) GetCertificate(ctx context.Context, req *domApi.GetCertificateRequest) (*projections.Certificate, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	for _, certificate := range ps {
		if certificate.ReferencedAggregateId == req.AggregateId && certificate.AggregateType == req.AggregateType {
			return certificate, nil
		}
	}

	return nil, esErrors.ErrProjectionNotFound
}
