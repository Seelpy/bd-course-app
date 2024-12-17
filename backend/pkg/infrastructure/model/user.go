package model

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
)

type User struct {
	ID       uuid.UUID
	AvatarID maybe.Maybe[uuid.UUID]
	Login    string
	Role     int
	Password string
	AboutMe  string
}
