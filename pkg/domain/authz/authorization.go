package authz

import (
	"context"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type authorizationCommandHandler struct {
}

func NewAuthorizationHandler() es.CommandHandler {
	return &authorizationCommandHandler{}
}

func (h *authorizationCommandHandler) HandleCommand(ctx context.Context, cmd es.Command) error {
	// TODO:
	// Get current users rolebindings from ctx
	// if !cmd.IsAuthorized(Admin, System, "") {
	// 	return fmt.Errorf("unauthorized")
	// }
	return nil
}
