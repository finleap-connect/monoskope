package commandhandler

import (
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commandhandler"
)

// apiServer is the implementation of the CommandHandler API
type apiServer struct {
	api.UnimplementedCommandHandlerServer
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer() (*apiServer, error) {
	s := &apiServer{}
	return s, nil
}
