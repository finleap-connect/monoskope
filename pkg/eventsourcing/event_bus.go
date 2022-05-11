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
)

// EventBusConnector can open and close connections.
type EventBusConnector interface {
	// Close closes the underlying connections
	Close() error
}

// EventBusPublisher publishes events on the underlying message bus.
type EventBusPublisher interface {
	EventBusConnector
	// PublishEvent publishes the event on the bus.
	PublishEvent(context.Context, Event) error
}

// EventBusConsumer notifies registered handlers on incoming events on the underlying message bus.
type EventBusConsumer interface {
	EventBusConnector
	// Matcher returns a new implementation specific matcher.
	Matcher() EventMatcher
	// AddHandler adds a handler for events matching one of the given EventMatcher.
	AddHandler(context.Context, EventHandler, ...EventMatcher) error
	// AddWorker behave similar to AddHandler but distributes events among the handlers with the same
	// work queue name according to the competing consumers pattern.
	AddWorker(context.Context, EventHandler, string, ...EventMatcher) error
}

// EventMatcher is an interface used to define what events should be consumed
type EventMatcher interface {
	// Any matches any event.
	Any() EventMatcher
	// MatchEventType matches a specific event type.
	MatchEventType(eventType EventType) EventMatcher
	// MatchAggregate matches a specific aggregate type.
	MatchAggregateType(aggregateType AggregateType) EventMatcher
}
