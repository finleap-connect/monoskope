package storage

import (
	"encoding/json"
)

type jsonEncoder struct{}

func (jsonEncoder) Marshal(data EventData) ([]byte, error) {
	if data != nil {
		return json.Marshal(data)
	}
	return nil, nil
}

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

func (jsonEncoder) String() string {
	return "json"
}
