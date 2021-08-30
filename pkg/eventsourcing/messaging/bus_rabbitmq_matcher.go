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

package messaging

import (
	"fmt"

	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// rabbitMatcher implements the EventMatcher interface for rabbitmq
type rabbitMatcher struct {
	routingKeyPrefix string
	eventType        string
	aggregateType    string
}

// Any matches any event.
func (m *rabbitMatcher) Any() evs.EventMatcher {
	m.eventType = "*"
	m.aggregateType = "*"
	return m
}

// MatchEventType matches a specific event type, nil events never match.
func (m *rabbitMatcher) MatchEventType(eventType evs.EventType) evs.EventMatcher {
	m.eventType = eventType.String()
	return m
}

// MatchAggregateType matches a specific aggregate type, nil events never match.
func (m *rabbitMatcher) MatchAggregateType(aggregateType evs.AggregateType) evs.EventMatcher {
	m.aggregateType = aggregateType.String()
	return m
}

// generateRoutingKey returns the routing key for events
func (m *rabbitMatcher) generateRoutingKey() string {
	return fmt.Sprintf("%s.%s.%s", m.routingKeyPrefix, m.aggregateType, m.eventType)
}
