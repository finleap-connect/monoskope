package storage

import (
	"encoding/json"
)

type jsonEncoder struct{}

// Marshal serializes EventData into JSON
func (jsonEncoder) Marshal(data EventData) ([]byte, error) {
	if data != nil {
		return json.Marshal(data)
	}
	return nil, nil
}

// Unmarshal deserializes JSON into EventData
func (jsonEncoder) Unmarshal(eventType EventType, raw []byte) (data EventData, err error) {
	if len(raw) == 0 {
		return nil, nil
	}

	if data, err = CreateEventData(eventType); err == nil {
		if err = json.Unmarshal(raw, data); err == nil {
			return data, nil
		}
	}
	return nil, err
}

// Returns the type of the encoder
func (jsonEncoder) String() string {
	return "json"
}
