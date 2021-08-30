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
	"context"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ctx = context.Background()
)

var _ = Describe("ServiceInformation", func() {
	It("can get service information", func() {
		conn, err := grpc.
			NewGrpcConnectionFactory(testEnv.GetApiAddr()).
			WithInsecure().
			WithBlock().
			Connect(ctx)
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()

		svc := api_common.NewServiceInformationServiceClient(conn)
		res, err := svc.GetServiceInformation(ctx, &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())

		var serviceInfos []*api_common.ServiceInformation
		for {
			// Read next
			serverInfo, err := res.Recv()

			// End of stream
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())

			// Append
			serviceInfos = append(serviceInfos, serverInfo)
		}
		Expect(serviceInfos).ToNot(BeNil())
		Expect(len(serviceInfos)).To(BeNumerically("==", 1))
		Expect(serviceInfos[0].GetName()).To(Equal(version.Name))
		Expect(serviceInfos[0].GetVersion()).To(Equal(version.Version))
		Expect(serviceInfos[0].GetCommit()).To(Equal(version.Commit))
	})
})
