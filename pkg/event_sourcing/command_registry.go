package event_sourcing

import (
	"context"
	"sync"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/anypb"
)

type CommandRegistry interface {
	CommandHandler
	RegisterCommand(func() Command) error
	CreateCommand(commandType CommandType, data *anypb.Any) (Command, error)
	SetHandler(handler CommandHandler, commandType CommandType) error
}

type commandRegistry struct {
	log      logger.Logger
	mutex    sync.RWMutex
	commands map[CommandType]func() Command
	handlers map[CommandType]CommandHandler
}

// newCommandRegistry creates a new command registry
func NewCommandRegistry() CommandRegistry {
	return &commandRegistry{
		log:      logger.WithName("command-registry"),
		commands: make(map[CommandType]func() Command),
		handlers: make(map[CommandType]CommandHandler),
	}
}

// RegisterCommand registers an command factory for a type. The factory is
// used to create concrete command types.
//
// An example would be:
//     RegisterCommand(func() Command { return &MyCommand{} })
func (r *commandRegistry) RegisterCommand(factory func() Command) error {
	cmd := factory()
	if cmd == nil {
		r.log.Info("factory does not create commands")
		return errors.ErrFactoryInvalid
	}

	commandType := cmd.CommandType()
	if commandType.String() == "" {
		r.log.Info("attempt to register empty command type")
		return errors.ErrEmptyCommandType
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.commands[commandType]; ok {
		r.log.Info("attempt to register command already registered", "commandType", commandType)
		return errors.ErrCommandTypeAlreadyRegistered
	}
	r.commands[commandType] = factory

	r.log.Info("command has been registered.", "commandType", commandType)

	return nil
}

// CreateCommand creates an command of a type with an ID using the factory
// registered with RegisterCommand.
func (r *commandRegistry) CreateCommand(commandType CommandType, data *anypb.Any) (Command, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if factory, ok := r.commands[commandType]; ok {
		cmd := factory()
		if err := cmd.SetData(data); err != nil {
			return nil, err
		}
		return cmd, nil
	}
	r.log.Info("trying to create a command of non-registered type", "commandType", commandType)
	return nil, errors.ErrCommandNotRegistered
}

// HandleCommand handles a command with a handler capable of handling it.
func (r *commandRegistry) HandleCommand(ctx context.Context, cmd Command) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if handler, ok := r.handlers[cmd.CommandType()]; ok {
		return handler.HandleCommand(ctx, cmd)
	}

	r.log.Info("trying to handle a command of non-registered type", "commandType", cmd.CommandType())
	return errors.ErrHandlerNotFound
}

// SetHandler adds a handler for a specific command.
func (r *commandRegistry) SetHandler(handler CommandHandler, commandType CommandType) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.handlers[commandType]; ok {
		r.log.Info("attempt to register command handler already registered", "commandType", commandType)
		return errors.ErrHandlerAlreadySet
	}

	r.handlers[commandType] = handler
	r.log.Info("command handler has been registered.", "commandType", commandType)

	return nil
}
