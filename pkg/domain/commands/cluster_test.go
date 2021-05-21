package commands

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

const shortDuration = 1 * time.Second

var _ = Describe("Unit Tests for Cluster Commands", func() {
	It("should read data for CreateClusterCommand from grpc message", func() {

		ctx, cancel := context.WithTimeout(context.Background(), shortDuration)
		defer cancel()

		agg := NewClusterAggregate(uuid.New())

		apiCommand, err := cmd.AddCommandData(
			cmd.CreateCommand(uuid.New(), commandTypes.CreateCluster),
			&cmdData.CreateCluster{
				Name:                "the one cluster",
				Label:               "one-cluster",
				ApiServerAddress:    "one.example.com",
				ClusterCACertBundle: []byte("This should be a certificate"),
			})
		Expect(err).ToNot(HaveOccurred())

		esCommand, err := s.cmdRegistry.CreateCommand(id, evs.CommandType(command.Type), command.Data)
		if err != nil {
			return nil, errors.TranslateToGrpcError(err)
		}

		err = agg.HandleCommand(ctx, esCommand)

		Expect(err).NotTo(HaveOccurred())

		Expect(ClusterAggregate(agg).GetName()).To(Equal("the one cluster"))

	})
})

func main() {
	// Pass a context with a timeout to tell a blocking function that it
	// should abandon its work after the timeout elapses.

}
