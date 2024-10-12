package model

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
)

type AuthorID = uuid.UUID

type Author struct {
	id         AuthorID
	avatarID   ImageID
	firstName  string
	secondName string
	middleName maybe.Maybe[string]
	nickname   maybe.Maybe[string]
}
