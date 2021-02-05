package event_sourcing

import "github.com/google/uuid"

type testProjection struct {
	Id      string
	Version uint64
}

func NewTestProjection(id uuid.UUID) *testProjection {
	return &testProjection{
		Id: id.String(),
	}
}

func (t testProjection) GetId() string {
	return t.Id
}

func (t testProjection) GetAggregateVersion() uint64 {
	return t.Version
}
