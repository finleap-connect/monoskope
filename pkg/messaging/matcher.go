package messaging

import "gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"

// EventMatcher
type EventMatcher interface {
	// MatchAny matches any event.
	MatchAny() EventMatcher
	// MatchEvent matches a specific event type, nil events never match.
	MatchEvent(eventType storage.EventType) EventMatcher
	// MatchAggregate matches a specific aggregate type, nil events never match.
	MatchAggregate(aggregateType storage.AggregateType) EventMatcher
}
