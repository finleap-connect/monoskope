package event_sourcing

import "github.com/google/uuid"

type testProjection struct {
	BaseProjection
}

func NewTestProjection() *testProjection {
	return &testProjection{
		BaseProjection: NewBaseProjection(uuid.New()),
	}
}
