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

	"github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	projectionsApi "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	es_repos "github.com/finleap-connect/monoskope/pkg/eventsourcing/repositories"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	adminRoleBinding.Role = common.Role_admin.String()
	adminRoleBinding.Scope = common.Scope_system.String()

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
