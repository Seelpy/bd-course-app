package query

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type BookQueryService interface {
	FindByID(bookID model.BookID) (BookOutput, error)
	List(page, size int) ([]BookOutput, error)
	CountBook(isPublished bool) (int, error)
}

type BookOutput struct {
	BookID      uuid.UUID
	Cover       maybe.Maybe[string]
	Title       string
	Description string
}

type bookQueryService struct {
	connection *sqlx.DB
}

func NewBookQueryService(connection *sqlx.DB) *bookQueryService {
	return &bookQueryService{connection: connection}
}

func (service *bookQueryService) FindByID(bookID model.BookID) (BookOutput, error) {
	const query = `
		SELECT b.book_id, i.path, b.title, b.description
		FROM book b
		LEFT OUTER JOIN image i ON b.cover_id = i.image_id
		WHERE b.book_id = ?;
	`

	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return BookOutput{}, err
	}

	var book sqlxBook
	err = service.connection.Get(&book, query, binaryBookID)
	if err != nil {
		return BookOutput{}, err
	}

	cover := maybe.Nothing[string]()
	if book.Cover.Valid {
		cover = maybe.Just(book.Cover.String)
	}

	return BookOutput{
		BookID:      book.BookID,
		Cover:       cover,
		Title:       book.Title,
		Description: book.Description,
	}, nil
}

func (service *bookQueryService) List(page, size int) ([]BookOutput, error) {
	const query = `
		SELECT b.book_id, i.path, b.title, b.description
		FROM book b
		LEFT OUTER JOIN image i ON b.cover_id = i.image_id
		WHERE b.is_publish = 1
		ORDER BY b.title
		LIMIT ? OFFSET ?;
	`

	offset := (page - 1) * size

	var sqlxBooks []sqlxBook
	err := service.connection.Select(&sqlxBooks, query, size, offset)
	if err != nil {
		return nil, err
	}

	bookOutputs := make([]BookOutput, len(sqlxBooks))
	for i, b := range sqlxBooks {
		cover := maybe.Nothing[string]()
		if b.Cover.Valid {
			cover = maybe.Just(b.Cover.String)
		}

		bookOutputs[i] = BookOutput{
			BookID:      b.BookID,
			Cover:       cover,
			Title:       b.Title,
			Description: b.Description,
		}
	}

	return bookOutputs, nil
}

func (service *bookQueryService) CountBook(isPublished bool) (int, error) {
	const query = `SELECT COUNT(*) FROM book b WHERE b.is_publish = ?`

	var countBook int
	err := service.connection.Get(&countBook, query, isPublished)
	if err != nil {
		return 0, err
	}

	return countBook, nil
}

type sqlxBook struct {
	BookID      uuid.UUID      `db:"book_id"`
	Cover       sql.NullString `db:"path"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
}
