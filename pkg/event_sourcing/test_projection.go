package event_sourcing

import "github.com/google/uuid"

type testProjection struct {
	id uuid.UUID
}

func NewTestProjection() *testProjection {
	return &testProjection{
		id: uuid.New(),
	}
}

func (p *testProjection) ID() uuid.UUID {
	return p.id
}
