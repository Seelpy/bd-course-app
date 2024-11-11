package model

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
)

type AuthorID = uuid.UUID

type Author struct {
	id         AuthorID
	avatarID   maybe.Maybe[ImageID]
	firstName  string
	secondName string
	middleName maybe.Maybe[string]
	nickname   maybe.Maybe[string]
}

func NewAuthor(
	id AuthorID,
	avatarID maybe.Maybe[ImageID],
	firstName string,
	secondName string,
	middleName maybe.Maybe[string],
	nickname maybe.Maybe[string],
) Author {
	return Author{
		id:         id,
		avatarID:   avatarID,
		firstName:  firstName,
		secondName: secondName,
		middleName: middleName,
		nickname:   nickname,
	}
}

func (author *Author) ID() AuthorID {
	return author.id
}

func (author *Author) AvatarID() maybe.Maybe[ImageID] {
	return author.avatarID
}

func (author *Author) FirstName() string {
	return author.firstName
}

func (author *Author) SecondName() string {
	return author.secondName
}

func (author *Author) MiddleName() maybe.Maybe[string] {
	return author.middleName
}

func (author *Author) Nickname() maybe.Maybe[string] {
	return author.nickname
}

func (author *Author) SetAvatarID(avatarID maybe.Maybe[ImageID]) {
	author.avatarID = avatarID
}

func (author *Author) SetFirstName(firstName string) {
	author.firstName = firstName
}

func (author *Author) SetSecondName(secondName string) {
	author.secondName = secondName
}

func (author *Author) SetMiddleName(middleName maybe.Maybe[string]) {
	author.middleName = middleName
}

func (author *Author) SetNickname(nickname maybe.Maybe[string]) {
	author.nickname = nickname
}
