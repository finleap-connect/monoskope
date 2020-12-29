package messaging

import "gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"

// EventMatcher is a func that can match event to a criteria.
type EventMatcher func(storage.Event) bool

// MatchAny matches any event.
func MatchAny() EventMatcher {
	return func(e storage.Event) bool {
		return true
	}
}

// MatchEvent matches a specific event type, nil events never match.
func MatchEvent(eventType storage.EventType) EventMatcher {
	return func(event storage.Event) bool {
		return event != nil && event.EventType() == eventType
	}
}

// MatchAggregate matches a specific aggregate type, nil events never match.
func MatchAggregate(aggregateType storage.AggregateType) EventMatcher {
	return func(event storage.Event) bool {
		return event != nil && event.AggregateType() == aggregateType
	}
}

// MatchAnyOf matches if any of several matchers matches.
func MatchAnyOf(matchers ...EventMatcher) EventMatcher {
	return func(event storage.Event) bool {
		for _, matcher := range matchers {
			if matcher(event) {
				return true
			}
		}
		return false
	}
}

// MatchAnyEventOf matches if any of several matchers matches.
func MatchAnyEventOf(eventTypes ...storage.EventType) EventMatcher {
	return func(event storage.Event) bool {
		for _, eventType := range eventTypes {
			if MatchEvent(eventType)(event) {
				return true
			}
		}
		return false
	}
}
