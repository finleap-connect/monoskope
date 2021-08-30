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
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("domain/certificate_repo", func() {

	var (
		expectedCert          = []byte("this is a certificate")
		expectedCACert        = []byte("this is the CA certificate")
		expectedAggregateType = "someaggregate"
	)

	certId := uuid.New()
	someAggregateId := uuid.New()
	userId := uuid.New()
	adminUser := &projections.User{User: &projectionsApi.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}}

	adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
	adminRoleBinding.UserId = adminUser.Id
	adminRoleBinding.Role = roles.Admin.String()
	adminRoleBinding.Scope = scopes.System.String()

	newCertificate := projections.NewCertificateProjection(certId).(*projections.Certificate)
	newCertificate.Certificate = &projectionsApi.Certificate{
		Id:                    certId.String(),
		ReferencedAggregateId: someAggregateId.String(),
		AggregateType:         expectedAggregateType,
		Certificate:           expectedCert,
		CaCertBundle:          expectedCACert,
	}
	newCertificate.Created = timestamp.New(time.Now())

	It("can retrieve the certificate", func() {
		inMemCertRepo := es_repos.NewInMemoryRepository()
		certRepo := NewCertificateRepository(inMemCertRepo)

		err := inMemCertRepo.Upsert(context.Background(), newCertificate)
		Expect(err).NotTo(HaveOccurred())
		cert, err := certRepo.GetCertificate(context.Background(),
			&domain.GetCertificateRequest{
				AggregateId:   someAggregateId.String(),
				AggregateType: expectedAggregateType,
			})
		Expect(err).NotTo(HaveOccurred())

		Expect(cert.Certificate.GetCertificate()).To(Equal(expectedCert))
		Expect(cert.Certificate.GetCaCertBundle()).To(Equal(expectedCACert))
	})

})
