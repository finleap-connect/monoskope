package messaging

import "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventstore/storage"

// EventMatcher is a func that can match event to a criteria.
type EventMatcher func(storage.Event) bool

// MatchAny matches any event.
func MatchAny() EventMatcher {
	return func(e storage.Event) bool {
		return true
	}
}

// MatchEvent matches a specific event type, nil events never match.
func MatchEvent(t storage.EventType) EventMatcher {
	return func(e storage.Event) bool {
		return e != nil && e.EventType() == t
	}
}

// MatchAggregate matches a specific aggregate type, nil events never match.
func MatchAggregate(t storage.AggregateType) EventMatcher {
	return func(e storage.Event) bool {
		return e != nil && e.AggregateType() == t
	}
}

// MatchAnyOf matches if any of several matchers matches.
func MatchAnyOf(matchers ...EventMatcher) EventMatcher {
	return func(e storage.Event) bool {
		for _, m := range matchers {
			if m(e) {
				return true
			}
		}
		return false
	}
}

// MatchAnyEventOf matches if any of several matchers matches.
func MatchAnyEventOf(types ...storage.EventType) EventMatcher {
	return func(e storage.Event) bool {
		for _, t := range types {
			if MatchEvent(t)(e) {
				return true
			}
		}
		return false
	}
}
