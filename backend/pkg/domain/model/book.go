package model

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
)

type BookID = uuid.UUID

type Book struct {
	id          BookID
	coverID     maybe.Maybe[ImageID]
	title       string
	description string
	isPublished bool
}

func NewBook(
	id BookID,
	coverID maybe.Maybe[ImageID],
	title string,
	description string,
	isPublished bool,
) Book {
	return Book{
		id:          id,
		coverID:     coverID,
		title:       title,
		description: description,
		isPublished: isPublished,
	}
}

func (book *Book) ID() BookID {
	return book.id
}

func (book *Book) CoverID() maybe.Maybe[ImageID] {
	return book.coverID
}

func (book *Book) Title() string {
	return book.title
}

func (book *Book) Description() string {
	return book.description
}

func (book *Book) IsPublished() bool {
	return book.isPublished
}

func (book *Book) SetCoverID(coverID maybe.Maybe[ImageID]) {
	book.coverID = coverID
}

func (book *Book) SetTitle(title string) {
	book.title = title
}

func (book *Book) SetDescription(description string) {
	book.description = description
}

func (book *Book) SetIsPublished(isPublished bool) {
	book.isPublished = isPublished
}
