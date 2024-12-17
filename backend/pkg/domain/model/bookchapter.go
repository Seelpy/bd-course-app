package model

import "github.com/gofrs/uuid"

type BookChapterID = uuid.UUID

type BookChapter struct {
	id     BookChapterID
	bookID BookID
	index  int
	title  string
}

func NewBookChapter(
	id BookChapterID,
	bookID BookID,
	index int,
	title string,
) BookChapter {
	return BookChapter{
		id:     id,
		bookID: bookID,
		index:  index,
		title:  title,
	}
}

func (bookChapter *BookChapter) ID() BookChapterID {
	return bookChapter.id
}

func (bookChapter *BookChapter) BookID() BookID {
	return bookChapter.bookID
}

func (bookChapter *BookChapter) Index() int {
	return bookChapter.index
}

func (bookChapter *BookChapter) Title() string {
	return bookChapter.title
}

func (bookChapter *BookChapter) SetTitle(title string) {
	bookChapter.title = title
}
