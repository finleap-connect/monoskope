package event_sourcing

import "encoding/json"

// EventData is any additional data for an event.
type EventData []byte

// ToEventData converts any given object into event data.
func ToEventData(data interface{}) (EventData, error) {
	return json.Marshal(data)
}
