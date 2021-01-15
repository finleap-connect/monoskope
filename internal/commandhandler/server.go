package commandhandler

import (
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commandhandler"
)

// commandHandler is the implementation of the CommandHandler API
type commandHandler struct {
	api.UnimplementedCommandHandlerServer
}

// NewCommandHandler returns a new configured instance of Server
func NewCommandHandler() (*commandHandler, error) {
	s := &commandHandler{}
	return s, nil
}
