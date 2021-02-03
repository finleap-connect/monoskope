package domain

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cmd_api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/user"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

var _ = Describe("domain", func() {
	It("can be set up", func() {
		err := es.Registry.RegisterCommand(func() es.Command { return &user.CreateUserCommand{} })
		Expect(err).NotTo(HaveOccurred())

		err = es.Registry.SetHandler(user.NewUserAggregate(uuid.New()), user.CreateUserType)
		Expect(err).NotTo(HaveOccurred())

		cmd := &cmd_api.CreateUserCommand{
			UserMetadata: &common.UserMetadata{
				Email: adminUser.Email,
				Name:  adminUser.Name,
			},
		}
		any := &anypb.Any{}
		Expect(any.MarshalFrom(cmd)).NotTo(HaveOccurred())
		esCmd, err := es.Registry.CreateCommand(user.CreateUserType, any)
		Expect(err).NotTo(HaveOccurred())

		err = es.Registry.HandleCommand(context.Background(), esCmd)
		Expect(err).NotTo(HaveOccurred())
	})
})
