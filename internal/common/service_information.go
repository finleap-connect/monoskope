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

package common

import (
	"github.com/finleap-connect/monoskope/internal/version"
	api_common "github.com/finleap-connect/monoskope/pkg/api/domain/common"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/metadata"
)

type serviceInformationService struct {
	api_common.UnimplementedServiceInformationServiceServer
	log logger.Logger
}

func (s *serviceInformationService) GetServiceInformation(e *empty.Empty, stream api_common.ServiceInformationService_GetServiceInformationServer) error {
	s.log.Info("Service information requested.")

	// Print headers
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		for k, vs := range md {
			for _, v := range vs {
				s.log.V(logger.DebugLevel).Info("Metadata provided.", "Key", k, "Value", v)
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
