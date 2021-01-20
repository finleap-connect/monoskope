package commands

import (
	"errors"
	"sync"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

// ErrFactoryInvalid is when a command factory creates nil commands.
var ErrFactoryInvalid = errors.New("factory does not create commands")

// ErrEmptyCommandType is when a command type given is empty.
var ErrEmptyCommandType = errors.New("command type must not be empty")

// ErrCommandTypeAlreadyRegistered is when a command was already registered.
var ErrCommandTypeAlreadyRegistered = errors.New("command type already registered")

// ErrCommandNotRegistered is when no command factory was registered.
var ErrCommandNotRegistered = errors.New("command not registered")

type CommandRegistry interface {
	RegisterCommand(factory func() Command) error
	UnregisterCommand(commandType CommandType) error
}

type commandRegistry struct {
	log      logger.Logger
	mutex    sync.RWMutex
	commands map[CommandType]func() Command
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() CommandRegistry {
	return &commandRegistry{
		log:      logger.WithName("command-registry"),
		commands: make(map[CommandType]func() Command),
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
		r.log.Info("registering duplicate types", "commandType", commandType)
		return ErrCommandTypeAlreadyRegistered
	}
	r.commands[commandType] = factory

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

	return nil
}

// CreateCommand creates an command of a type with an ID using the factory
// registered with RegisterCommand.
func (r *commandRegistry) CreateCommand(commandType CommandType) (Command, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if factory, ok := r.commands[commandType]; ok {
		return factory(), nil
	}
	return nil, ErrCommandNotRegistered
}
