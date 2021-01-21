package events

import (
	"encoding/json"

	"github.com/google/uuid"
	. "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type UserRoleAddedData struct {
	UserId  uuid.UUID `json:",omitempty"`
	Role    string    `json:",omitempty"`
	Context string    `json:",omitempty"`
}

func UserRoleAddedDataFromEventData(ed EventData) (*UserRoleAddedData, error) {
	data := &UserRoleAddedData{}
	if err := json.Unmarshal(ed, data); err != nil {
		return nil, err
	}
	return data, nil
}
