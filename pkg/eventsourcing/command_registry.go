// Copyright 2022 Monoskope Authors
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

package eventsourcing

import (
	"context"
	"sync"

	"github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

type CommandRegistry interface {
	CommandHandler
	RegisterCommand(func(uuid.UUID) Command)
	CreateCommand(id uuid.UUID, commandType CommandType, data *anypb.Any) (Command, error)
	GetRegisteredCommandTypes() []CommandType
	SetHandler(handler CommandHandler, commandType CommandType)
}

type commandRegistry struct {
	log      logger.Logger
	mutex    sync.RWMutex
	commands map[CommandType]func(uuid.UUID) Command
	handlers map[CommandType]CommandHandler
}

var DefaultCommandRegistry CommandRegistry

func init() {
	DefaultCommandRegistry = NewCommandRegistry()
}

// newCommandRegistry creates a new command registry
func NewCommandRegistry() CommandRegistry {
	return &commandRegistry{
		log:      logger.WithName("command-registry"),
		commands: make(map[CommandType]func(uuid.UUID) Command),
		handlers: make(map[CommandType]CommandHandler),
	}
}

// GetRegisteredCommandTypes returns a list with all registered command types.
func (r *commandRegistry) GetRegisteredCommandTypes() []CommandType {
	keys := make([]CommandType, 0, len(r.commands))
	for k := range r.commands {
		keys = append(keys, k)
	}
	return keys
}

// RegisterCommand registers an command factory for a type. The factory is
// used to create concrete command types.
//
// An example would be:
//
//	RegisterCommand(func() Command { return &MyCommand{} })
func (r *commandRegistry) RegisterCommand(factory func(uuid.UUID) Command) {
	cmd := factory(uuid.Nil)
	if cmd == nil {
		r.log.Info("factory does not create commands")
		panic(errors.ErrFactoryInvalid)
	}

	commandType := cmd.CommandType()
	if commandType.String() == "" {
		r.log.Info("attempt to register empty command type")
		panic(errors.ErrEmptyCommandType)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.commands[commandType]; ok {
		r.log.Info("attempt to register command already registered", "commandType", commandType)
		panic(errors.ErrCommandTypeAlreadyRegistered)
	}
	r.commands[commandType] = factory

	r.log.V(logger.DebugLevel).Info("command has been registered.", "commandType", commandType)
}

// CreateCommand creates an command of a type with an ID using the factory
// registered with RegisterCommand.
func (r *commandRegistry) CreateCommand(id uuid.UUID, commandType CommandType, data *anypb.Any) (Command, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if factory, ok := r.commands[commandType]; ok {
		cmd := factory(id)
		if data != nil {
			if err := cmd.SetData(data); err != nil {
				return nil, err
			}
		}
		return cmd, nil
	}
	r.log.Info("trying to create a command of non-registered type", "commandType", commandType)
	return nil, errors.ErrCommandNotRegistered
}

// HandleCommand handles a command with a handler capable of handling it.
func (r *commandRegistry) HandleCommand(ctx context.Context, cmd Command) (*CommandReply, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if handler, ok := r.handlers[cmd.CommandType()]; ok {
		return handler.HandleCommand(ctx, cmd)
	}

	r.log.Info("trying to handle a command of non-registered type", "commandType", cmd.CommandType())
	return nil, errors.ErrHandlerNotFound
}

// SetHandler adds a handler for a specific command.
func (r *commandRegistry) SetHandler(handler CommandHandler, commandType CommandType) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.handlers[commandType]; ok {
		r.log.Info("attempt to register command handler already registered", "commandType", commandType)
		panic(errors.ErrHandlerAlreadySet)
	}

	r.handlers[commandType] = handler
	r.log.V(logger.DebugLevel).Info("command handler has been registered.", "commandType", commandType)
}
