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

func (jsonEncoder) Unmarshal(raw []byte, data interface{}) error {
	if len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, data); err != nil {
		return err
	}
	return nil
}

func (jsonEncoder) String() string {
	return "json"
}
