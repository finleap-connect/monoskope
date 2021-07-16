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

var (
	expectedCert          = []byte("this is a certificate")
	expectedCACert        = []byte("this is the CA certificate")
	expectedAggregateType = "certificate"
)

var _ = Describe("domain/certificate_repo", func() {

	certId := uuid.New()

	userId := uuid.New()
	adminUser := &projections.User{User: &projectionsApi.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}}

	adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
	adminRoleBinding.UserId = adminUser.Id
	adminRoleBinding.Role = roles.Admin.String()
	adminRoleBinding.Scope = scopes.System.String()

	newCertificate := projections.NewCertificateProjection(certId).(*projections.Certificate)
	newCertificate.Certificate = &projectionsApi.Certificate{
		ReferencedAggregateId: certId.String(),
		AggregateType:         expectedAggregateType,
		Certificate:           expectedCert,
		CaCertBundle:          expectedCACert,
	}
	newCertificate.Created = timestamp.New(time.Now())

	It("can retrieve the certificate", func() {
		inMemoryRoleRepo := es_repos.NewInMemoryRepository()
		err := inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)
		Expect(err).NotTo(HaveOccurred())

		userRoleBindingRepo := NewUserRoleBindingRepository(inMemoryRoleRepo)

		inMemoryUserRepo := es_repos.NewInMemoryRepository()
		userRepo := NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)
		err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())

		inMemCertRepo := es_repos.NewInMemoryRepository()
		certRepo := NewCertificateRepository(inMemCertRepo, userRepo)

		err = inMemCertRepo.Upsert(context.Background(), newCertificate)
		Expect(err).NotTo(HaveOccurred())
		cert, err := certRepo.GetCertificate(context.Background(),
			&domain.GetCertificateRequest{
				AggregateId:   certId.String(),
				AggregateType: expectedAggregateType,
			})
		Expect(err).NotTo(HaveOccurred())

		Expect(cert.Certificate.GetCertificate()).To(Equal(expectedCert))
		Expect(cert.Certificate.GetCaCertBundle()).To(Equal(expectedCACert))
	})

})
