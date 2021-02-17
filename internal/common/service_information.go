package common

import (
	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
)

type serviceInformationService struct {
	api_common.UnimplementedServiceInformationServiceServer
}

func (s *serviceInformationService) GetServiceInformation(e *empty.Empty, stream api_common.ServiceInformationService_GetServiceInformationServer) error {
	err := stream.Send(&api_common.ServiceInformation{
		Name:    version.Name,
		Version: version.Version,
		Commit:  version.Commit,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewServiceInformationService() *serviceInformationService {
	return &serviceInformationService{}
}
