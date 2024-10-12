package model

import "time"

type ReadingSession struct {
	bookID        BookID
	bookChapterID BookChapterID
	userID        UserID
	lastReadTime  time.Time
}
