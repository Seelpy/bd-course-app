package service

import (
	"server/pkg/domain/model"
	"time"
)

type ReadingSessionService interface {
	StoreReadingSession(input StoreReadingSessionInput) error
}

type readingSessionService struct {
	readingSessionRepo ReadingSessionRepository
}

func NewReadingSessionService(
	readingSessionRepo ReadingSessionRepository,
) *readingSessionService {
	return &readingSessionService{
		readingSessionRepo: readingSessionRepo,
	}
}

type ReadingSessionRepository interface {
	Store(readingSession model.ReadingSession) error
}

type StoreReadingSessionInput struct {
	BookID        model.BookID
	BookChapterID model.BookChapterID
	UserID        model.UserID
}

func (service *readingSessionService) StoreReadingSession(input StoreReadingSessionInput) error {
	readingSession := model.NewReadingSession(
		input.BookID,
		input.BookChapterID,
		input.UserID,
		time.Now(),
	)

	return service.readingSessionRepo.Store(readingSession)
}
