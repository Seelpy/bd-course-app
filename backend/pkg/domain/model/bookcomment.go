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
