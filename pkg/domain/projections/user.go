package projections

import "github.com/google/uuid"

type User struct {
	Email string
}

func (u *User) ID() uuid.UUID {
	return uuid.Nil
}
