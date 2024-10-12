package model

import (
	"github.com/gofrs/uuid"
	"time"
)

type BookCommentID = uuid.UUID

type BookComment struct {
	id        BookCommentID
	bookID    BookID
	userID    UserID
	comment   string
	createdAt time.Time
}

func NewBookComment(
	id BookCommentID,
	bookID BookID,
	userID UserID,
	comment string,
	createdAt time.Time,
) BookComment {
	return BookComment{
		id:        id,
		bookID:    bookID,
		userID:    userID,
		comment:   comment,
		createdAt: createdAt,
	}
}

func (bookComment *BookComment) ID() BookCommentID {
	return bookComment.id
}

func (bookComment *BookComment) BookID() BookID {
	return bookComment.bookID
}

func (bookComment *BookComment) UserID() UserID {
	return bookComment.userID
}

func (bookComment *BookComment) Comment() string {
	return bookComment.comment
}

func (bookComment *BookComment) CreatedAt() time.Time {
	return bookComment.createdAt
}

func (bookComment *BookComment) SetComment(comment string) {
	bookComment.comment = comment
}
