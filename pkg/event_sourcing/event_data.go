package event_sourcing

import "encoding/json"

// EventData is any additional data for an event.
type EventData []byte

// MarshalEventData marshals any given object into event data.
func MarshalEventData(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// UnmarshalEventData unmarshals event data to a given type.
func UnmarshalEventData(eventData []byte, v interface{}) error {
	return json.Unmarshal(eventData, v)
}
