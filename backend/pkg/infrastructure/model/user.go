package model

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
)

type User struct {
	ID       uuid.UUID
	Avatar   maybe.Maybe[string]
	Login    string
	Role     int
	Password string
	AboutMe  string
}
