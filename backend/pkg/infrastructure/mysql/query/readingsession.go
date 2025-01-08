package query

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type ReadingSessionQueryService interface {
	GetLastReadingSession(bookID model.BookID, userID model.UserID) (ReadingSessionOutput, error)
}

type ReadingSessionOutput struct {
	BookChapterId uuid.UUID
}

type readingSessionQueryService struct {
	connection *sqlx.DB
}

func NewReadingSessionQueryService(connection *sqlx.DB) *readingSessionQueryService {
	return &readingSessionQueryService{connection: connection}
}

func (service *readingSessionQueryService) GetLastReadingSession(
	bookID model.BookID,
	userID model.UserID,
) (ReadingSessionOutput, error) {
	const query = `
		SELECT rs.book_chapter_id
		FROM reading_session rs
		WHERE rs.book_id = ? AND rs.user_id = ?
		ORDER BY rs.last_read_time DESC
	`

	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return ReadingSessionOutput{}, err
	}
	binaryUserID, err := uuid.UUID(userID).MarshalBinary()
	if err != nil {
		return ReadingSessionOutput{}, err
	}

	var sqlxReadingSessions []sqlxReadingSession
	err = service.connection.Select(&sqlxReadingSessions, query, binaryBookID, binaryUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return ReadingSessionOutput{}, model.ErrSessionReadingNotFound
	}
	if err != nil {
		return ReadingSessionOutput{}, err
	}

	return ReadingSessionOutput{
		BookChapterId: sqlxReadingSessions[0].BookChapterId,
	}, nil
}

type sqlxReadingSession struct {
	BookChapterId uuid.UUID `db:"book_chapter_id"`
}
