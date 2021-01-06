package messaging

import (
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

// RabbitMatcher implements the EventMatcher interface for rabbitmq
type RabbitMatcher struct {
	routingKeyPrefix string
	eventType        string
	aggregateType    string
}

// Any matches any event.
func (m *RabbitMatcher) Any() EventMatcher {
	m.eventType = "*"
	m.aggregateType = "*"
	return m
}

// MatchEventType matches a specific event type, nil events never match.
func (m *RabbitMatcher) MatchEventType(eventType storage.EventType) EventMatcher {
	m.eventType = string(eventType)
	return m
}

// MatchAggregateType matches a specific aggregate type, nil events never match.
func (m *RabbitMatcher) MatchAggregateType(aggregateType storage.AggregateType) EventMatcher {
	m.aggregateType = string(aggregateType)
	return m
}

// generateRoutingKey returns the routing key for events
func (m *RabbitMatcher) generateRoutingKey() string {
	return fmt.Sprintf("%s.%s.%s", m.routingKeyPrefix, m.aggregateType, m.eventType)
}
