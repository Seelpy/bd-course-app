package query

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type BookChapterQueryService interface {
	ListByBookID(bookID model.BookID) ([]BookChapterOutput, error)
}

type BookChapterOutput struct {
	BookChapterID uuid.UUID
	Index         int
	Title         string
}

type bookChapterQueryService struct {
	connection *sqlx.DB
}

func NewBookChapterQueryService(connection *sqlx.DB) *bookChapterQueryService {
	return &bookChapterQueryService{
		connection: connection,
	}
}

func (service *bookChapterQueryService) ListByBookID(bookID model.BookID) ([]BookChapterOutput, error) {
	const query = `
		SELECT bc.book_chapter_id, bc.chapter_index, bc.title
		FROM book_chapter bc
		WHERE bc.book_id = ?
		ORDER BY bc.book_chapter_id ASC
	`

	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return nil, err
	}

	var sqlxBookChapters []sqlxBookChapter
	err = service.connection.Select(&sqlxBookChapters, query, binaryBookID)
	if err != nil {
		return nil, err
	}

	bookChapterOutputs := make([]BookChapterOutput, len(sqlxBookChapters))
	for i, b := range sqlxBookChapters {
		bookChapterOutputs[i] = BookChapterOutput{
			BookChapterID: b.BookChapterID,
			Index:         b.Index,
			Title:         b.Title,
		}
	}

	return bookChapterOutputs, nil
}

type sqlxBookChapter struct {
	BookChapterID uuid.UUID `db:"book_chapter_id"`
	Index         int       `db:"chapter_index"`
	Title         string    `db:"title"`
}
