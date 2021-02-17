package repositories

import "github.com/google/uuid"

type testProjection struct {
	id      uuid.UUID
	version uint64
}

func newTestProjection(id uuid.UUID) *testProjection {
	return &testProjection{
		id: id,
	}
}

func (t testProjection) ID() uuid.UUID {
	return t.id
}

func (t testProjection) Version() uint64 {
	return t.version
}

func (t testProjection) IncrementVersion() {
	t.version++
}
