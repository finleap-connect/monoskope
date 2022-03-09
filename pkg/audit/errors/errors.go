// Copyright 2021 Monoskope Authors
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

// Event formatter Errors
var (
	// ErrMissingFormatterImplementationForEventType is when the event type is not supported by the formatter
	ErrMissingFormatterImplementationForEventType = errors.New("event type is not supported by formatter")

	// ErrEventFormatterInvalid is when an event formatter is invalid e.g. nil.
	ErrEventFormatterInvalid = errors.New("invalid event formatter")

	// ErrEmptyEventType is when an event type given is empty.
	ErrEmptyEventType = errors.New("event type must not be empty")

	// ErrEventFormatterForEventTypeAlreadyRegistered is when an event formatter for event type was already registered.
	ErrEventFormatterForEventTypeAlreadyRegistered = errors.New("event formatter for event type already registered")

	// ErrEventFormatterForEventTypeNotRegistered is when no event formatter was registered for the event type.
	ErrEventFormatterForEventTypeNotRegistered = errors.New("event formatter for event type not registered")
)
