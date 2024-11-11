package model

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
)

type UserID uuid.UUID

type UserRole int16

const (
	Admin UserRole = iota
	Client
)

type User struct {
	id       UserID
	avatarID maybe.Maybe[ImageID]
	login    string
	role     UserRole
	password string
	aboutMe  string
}

func NewUser(
	id UserID,
	avatarID maybe.Maybe[ImageID],
	login string,
	role UserRole,
	password string,
	aboutMe string,
) User {
	return User{
		id:       id,
		avatarID: avatarID,
		login:    login,
		role:     role,
		password: password,
		aboutMe:  aboutMe,
	}
}

func (user *User) ID() UserID {
	return user.id
}

func (user *User) AvatarID() maybe.Maybe[ImageID] {
	return user.avatarID
}

func (user *User) Login() string {
	return user.login
}

func (user *User) Role() UserRole {
	return user.role
}

func (user *User) Password() string {
	return user.password
}

func (user *User) AboutMe() string {
	return user.aboutMe
}

func (user *User) SetAvatarID(avatarID maybe.Maybe[ImageID]) {
	user.avatarID = avatarID
}

func (user *User) SetLogin(login string) {
	user.login = login
}

func (user *User) SetPassword(password string) {
	user.password = password
}

func (user *User) SetAboutMe(aboutMe string) {
	user.aboutMe = aboutMe
}
