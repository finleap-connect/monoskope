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
