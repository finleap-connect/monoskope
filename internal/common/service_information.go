package common

import (
	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"google.golang.org/grpc/metadata"
)

type serviceInformationService struct {
	api_common.UnimplementedServiceInformationServiceServer
	log logger.Logger
}

func (s *serviceInformationService) GetServiceInformation(e *empty.Empty, stream api_common.ServiceInformationService_GetServiceInformationServer) error {
	s.log.Info("Service information requested.")

	if util.GetOperationMode() == util.DEVELOPMENT {
		// Print headers
		if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
			for k, vs := range md {
				for _, v := range vs {
					s.log.Info("Metadata provided.", "Key", k, "Value", v)
				}
			}
		}
	}

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
	return &serviceInformationService{
		log: logger.WithName("ServiceInformationService"),
	}
}
