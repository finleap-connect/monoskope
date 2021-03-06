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

package errors

import "errors"

// EventStorage Errors
var (
	// ErrNoEventsToAppend is when no events are available to append.
	ErrNoEventsToAppend = errors.New("no events to append")

	// ErrIncorrectEventAggregateVersion is when an event is for an other version of the aggregate.
	ErrIncorrectAggregateVersion = errors.New("mismatching event aggregate version")

	// ErrAggregateVersionAlreadyExists is when an event is referencing an older version of the aggregate than is stored.
	ErrAggregateVersionAlreadyExists = errors.New("event aggregate version already exists in store")

	// ErrInvalidAggregateType is when an event is for a different type of aggregate.
	ErrInvalidAggregateType = errors.New("mismatching event aggregate type")

	// ErrCouldNotSaveEvents is when events could not be saved.
	ErrCouldNotSaveEvents = errors.New("could not save events")

	// ErrCouldNotConnect is when store could not connect to the underlying storage
	ErrCouldNotConnect = errors.New("could not connect to storage")

	// ErrConnectionClosed is when connection with underlying storage has been closed
	ErrConnectionClosed = errors.New("connection to storage closed")
)

// MessageBus Errors
var (
	// ErrCouldNotMarshalEvent is when an event could not be marshaled.
	ErrCouldNotMarshalEvent = errors.New("could not marshal event")

	// ErrCouldNotPublishEvent is when cannot send event to message bus
	ErrCouldNotPublishEvent = errors.New("could not publish event")

	// ErrMatcherMustNotBeNil is when an empty matcher has been provided
	ErrMatcherMustNotBeNil = errors.New("matcher must not be nil")

	// ErrHandlerMustNotBeNil is when an empty handler has been provided
	ErrHandlerMustNotBeNil = errors.New("handler must not be nil")

	// ErrMessageNotConnected is when there is no connection
	ErrMessageNotConnected = errors.New("message bus not connected")

	// ErrMessageBusConnection is when an unexpected error on message bus occurred
	ErrMessageBusConnection = errors.New("unexpected error on message bus occurred")

	// ErrCouldNotAddHandler is when an handler could not be added
	ErrCouldNotAddHandler = errors.New("could not add handler")

	// ErrContextDeadlineExceeded is when execution has been aborted since the context deadline has been exceeded
	ErrContextDeadlineExceeded = errors.New("context deadline exceeded")

	// ErrCouldNotParseAggregateId is when an aggregate id could not be parsed as uuid
	ErrCouldNotParseAggregateId = errors.New("could not parse aggregate id")

	// ErrConfigNameRequired is when the config doesn't include a name.
	ErrConfigNameRequired = errors.New("name must not be empty")

	// ErrConfigUrlRequired is when the config doesn't include a name.
	ErrConfigUrlRequired = errors.New("url must not be empty")
)

// Command Registry Errors
var (
	// ErrFactoryInvalid is when a command factory creates nil commands.
	ErrFactoryInvalid = errors.New("factory does not create commands")

	// ErrEmptyCommandType is when a command type given is empty.
	ErrEmptyCommandType = errors.New("command type must not be empty")

	// ErrCommandTypeAlreadyRegistered is when a command was already registered.
	ErrCommandTypeAlreadyRegistered = errors.New("command type already registered")

	// ErrCommandNotRegistered is when no command factory was registered.
	ErrCommandNotRegistered = errors.New("command not registered")

	// ErrHandlerAlreadySet is when a handler is already registered for a command.
	ErrHandlerAlreadySet = errors.New("handler is already set")

	// ErrHandlerNotFound is when no handler can be found.
	ErrHandlerNotFound = errors.New("no handlers for command")
)

// Repository Errors
var (
	// ErrProjectionNotFound is when the requested Projection was not found in the repository.
	ErrProjectionNotFound = errors.New("projection not found")
)

var (
	// ErrInvalidEventType is when a projection received an event for a different type.
	ErrInvalidEventType = errors.New("mismatching event type")

	// ErrInvalidProjectionType is when a projection is invalid.
	ErrInvalidProjectionType = errors.New("mismatching projection type")
)

// Aggregate Registry Errors
var (
	// ErrEmptyAggregateType is when a aggregate type given is empty.
	ErrEmptyAggregateType = errors.New("aggregate type must not be empty")

	// ErrAggregateTypeAlreadyRegistered is when a aggregate was already registered.
	ErrAggregateTypeAlreadyRegistered = errors.New("aggregate type already registered")

	// ErrAggregateNotRegistered is when no aggregate factory was registered.
	ErrAggregateNotRegistered = errors.New("aggregate not registered")
)
