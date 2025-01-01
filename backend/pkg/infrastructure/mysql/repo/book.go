package repo

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type BookRepository struct {
	connection *sqlx.DB
}

func NewBookRepository(connection *sqlx.DB) *BookRepository {
	return &BookRepository{connection}
}

func (repo *BookRepository) NextID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (repo *BookRepository) Store(book model.Book) error {
	const query = `
		INSERT INTO
			book (
			      book_id,
			      title,
			      description,
			      is_publish
			)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			title = VALUES(title),
			description = VALUES(description),
			is_publish = VALUES(is_publish)
	`

	binaryBookID, err := uuid.UUID(book.ID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query,
		binaryBookID,
		book.Title(),
		book.Description(),
		book.IsPublished(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *BookRepository) Delete(bookID model.BookID) error {
	const query = `DELETE FROM book WHERE book_id = ?`

	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryBookID)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	return err
}

func (repo *BookRepository) FindByID(bookID model.BookID) (model.Book, error) {
	const query = `
		SELECT
			book_id,
			description,
			title,
			is_publish
		FROM book
		WHERE book_id = ?
`

	var book sqlxBook
	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return model.Book{}, err
	}

	err = repo.connection.Get(&book, query, binaryBookID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Book{}, model.ErrUserNotFound
	}
	if err != nil {
		return model.Book{}, err
	}

	return model.NewBook(
		model.BookID(bookID),
		maybe.Nothing[model.ImageID](),
		book.Title,
		book.Description,
		book.IsPublished,
	), nil
}

type sqlxBook struct {
	Description string `db:"description"`
	Title       string `db:"title"`
	IsPublished bool   `db:"is_publish"`
}
