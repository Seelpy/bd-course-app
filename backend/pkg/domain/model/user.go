package model

import "github.com/gofrs/uuid"

type UserID uuid.UUID

type User struct {
	id       UserID
	avatarID ImageID
	login    string
	role     int16
	password string
	aboutMe  string
}
