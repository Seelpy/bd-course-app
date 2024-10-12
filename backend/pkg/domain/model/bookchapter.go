package model

import "github.com/gofrs/uuid"

type BookChapterID = uuid.UUID

type BookChapter struct {
	id     BookChapterID
	bookID BookID
	index  int
	title  string
}
