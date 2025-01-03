package repo

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type readingSessionRepository struct {
	connection *sqlx.DB
}

func NewReadingSessionRepository(connection *sqlx.DB) *readingSessionRepository {
	return &readingSessionRepository{connection: connection}
}

func (repo *readingSessionRepository) Store(readingSession model.ReadingSession) error {
	const query = `
		INSERT INTO
			reading_session (
			      book_id,
				  book_chapter_id,
				  user_id,
				  last_read_time
			)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			last_read_time = VALUES(last_read_time)
	`

	binaryBookID, err := uuid.UUID(readingSession.BookID()).MarshalBinary()
	if err != nil {
		return err
	}
	binaryBookChapterID, err := uuid.UUID(readingSession.BookChapterID()).MarshalBinary()
	if err != nil {
		return err
	}
	binaryUserID, err := uuid.UUID(readingSession.UserID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query,
		binaryBookID,
		binaryBookChapterID,
		binaryUserID,
		readingSession.LastReadTime(),
	)

	return err
}
