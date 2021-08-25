package queryhandler

import (
	"context"
	"time"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/grpc"
)

// certificateServer is the implementation of the CertificateService API
type certificateServer struct {
	api.UnimplementedCertificateServer

	repoCertificate repositories.ReadOnlyCertificateRepository
}

// NewCertificateServiceServer returns a new configured instance of certificateServiceServer
func NewCertificateServer(certificateRepo repositories.ReadOnlyCertificateRepository) *certificateServer {
	return &certificateServer{
		repoCertificate: certificateRepo,
	}
}

func NewCertificateClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, api.CertificateClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(queryHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, api.NewCertificateClient(conn), nil
}

// GetById returns the certificate found by the given id.
func (s *certificateServer) GetCertificate(ctx context.Context, gcreq *api.GetCertificateRequest) (*projections.Certificate, error) {
	certificate, err := s.repoCertificate.GetCertificate(ctx, gcreq)
	if err != nil {
		return nil, err
	}
	return certificate.Proto(), nil
}
