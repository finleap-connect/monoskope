package domain

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cmd_api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	types "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

var _ = Describe("domain", func() {
	adminUser := &projections.User{Email: "admin@monoskope.io", Name: "admin"}

	It("can be set up", func() {
		registry := es.NewCommandRegistry()

		err := registry.RegisterCommand(func() es.Command { return &commands.CreateUserCommand{} })
		Expect(err).NotTo(HaveOccurred())

		err = registry.SetHandler(aggregates.NewUserAggregate(uuid.New()), types.CreateUserType)
		Expect(err).NotTo(HaveOccurred())

		cmd := &cmd_api.CreateUserCommand{
			UserMetadata: &common.UserMetadata{
				Email: adminUser.Email,
				Name:  adminUser.Name,
			},
		}
		any := &anypb.Any{}
		Expect(any.MarshalFrom(cmd)).NotTo(HaveOccurred())
		esCmd, err := registry.CreateCommand(types.CreateUserType, any)
		Expect(err).NotTo(HaveOccurred())

		err = registry.HandleCommand(context.Background(), esCmd)
		Expect(err).NotTo(HaveOccurred())
	})
})
