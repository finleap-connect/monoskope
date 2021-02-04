package event_sourcing

import (
	"context"
	"errors"
	"sync"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/anypb"
)

// ErrFactoryInvalid is when a command factory creates nil commands.
var ErrFactoryInvalid = errors.New("factory does not create commands")

// ErrEmptyCommandType is when a command type given is empty.
var ErrEmptyCommandType = errors.New("command type must not be empty")

// ErrCommandTypeAlreadyRegistered is when a command was already registered.
var ErrCommandTypeAlreadyRegistered = errors.New("command type already registered")

// ErrCommandNotRegistered is when no command factory was registered.
var ErrCommandNotRegistered = errors.New("command not registered")

// ErrHandlerAlreadySet is when a handler is already registered for a command.
var ErrHandlerAlreadySet = errors.New("handler is already set")

// ErrHandlerNotFound is when no handler can be found.
var ErrHandlerNotFound = errors.New("no handlers for command")

type CommandRegistry interface {
	CommandHandler
	RegisterCommand(func() Command) error
	UnregisterCommand(commandType CommandType) error
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
		return ErrFactoryInvalid
	}

	commandType := cmd.CommandType()
	if commandType.String() == "" {
		r.log.Info("attempt to register empty command type")
		return ErrEmptyCommandType
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.commands[commandType]; ok {
		r.log.Info("attempt to register command already registered", "commandType", commandType)
		return ErrCommandTypeAlreadyRegistered
	}
	r.commands[commandType] = factory

	r.log.Info("command has been registered.", "commandType", commandType)

	return nil
}

// UnregisterCommand removes the registration of the command factory for
// a type. This is mainly useful in mainenance situations where the command type
// needs to be switched at runtime.
func (r *commandRegistry) UnregisterCommand(commandType CommandType) error {
	if commandType == CommandType("") {
		r.log.Info("attempt to register empty command type")
		return ErrEmptyCommandType
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, ok := r.commands[commandType]; !ok {
		r.log.Info("unregister of non-registered type", "commandType", commandType)
		return ErrCommandNotRegistered
	}
	delete(r.commands, commandType)

	r.log.Info("command has been unregistered.", "commandType", commandType)

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
	return nil, ErrCommandNotRegistered
}

// HandleCommand handles a command with a handler capable of handling it.
func (r *commandRegistry) HandleCommand(ctx context.Context, cmd Command) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if handler, ok := r.handlers[cmd.CommandType()]; ok {
		return handler.HandleCommand(ctx, cmd)
	}

	r.log.Info("trying to handle a command of non-registered type", "commandType", cmd.CommandType())
	return ErrHandlerNotFound
}

// SetHandler adds a handler for a specific command.
func (r *commandRegistry) SetHandler(handler CommandHandler, commandType CommandType) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.handlers[commandType]; ok {
		r.log.Info("attempt to register command handler already registered", "commandType", commandType)
		return ErrHandlerAlreadySet
	}

	r.handlers[commandType] = handler
	r.log.Info("command handler has been registered.", "commandType", commandType)

	return nil
}
