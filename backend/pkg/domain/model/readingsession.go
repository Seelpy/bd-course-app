package model

import "time"

type ReadingSession struct {
	bookID        BookID
	bookChapterID BookChapterID
	userID        UserID
	lastReadTime  time.Time
}

func NewReadingSession(
	bookID BookID,
	bookChapterID BookChapterID,
	userID UserID,
	lastReadTime time.Time,
) ReadingSession {
	return ReadingSession{
		bookID:        bookID,
		bookChapterID: bookChapterID,
		userID:        userID,
		lastReadTime:  lastReadTime,
	}
}

func (readingSession *ReadingSession) BookID() BookID {
	return readingSession.bookID
}

func (readingSession *ReadingSession) BookChapterID() BookChapterID {
	return readingSession.bookChapterID
}

func (readingSession *ReadingSession) UserID() UserID {
	return readingSession.userID
}

func (readingSession *ReadingSession) LastReadTime() time.Time {
	return readingSession.lastReadTime
}

func (readingSession *ReadingSession) SetLastReadTime(lastReadTime time.Time) {
	readingSession.lastReadTime = lastReadTime
}
