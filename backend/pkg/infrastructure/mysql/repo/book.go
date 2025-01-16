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
			      is_publish,
			      cover_id
			)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			title = VALUES(title),
			description = VALUES(description),
			is_publish = VALUES(is_publish),
			cover_id = VALUES(cover_id)
	`

	binaryBookID, err := uuid.UUID(book.ID()).MarshalBinary()
	if err != nil {
		return err
	}

	var coverID *[]byte
	if imageID, ok := book.CoverID().Get(); ok {
		uid, err2 := uuid.UUID(imageID).MarshalBinary()
		if err2 != nil {
			return err2
		}

		coverID = &uid
	} else {
		coverID = nil
	}

	_, err = repo.connection.Exec(query,
		binaryBookID,
		book.Title(),
		book.Description(),
		book.IsPublished(),
		coverID,
	)

	return err
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

	count, err := result.RowsAffected()
	if count == 0 {
		return model.ErrBookNotFound
	}

	return err
}

func (repo *BookRepository) FindByID(bookID model.BookID) (model.Book, error) {
	const query = `
		SELECT
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
		return model.Book{}, model.ErrBookNotFound
	}
	if err != nil {
		return model.Book{}, err
	}

	return model.NewBook(
		bookID,
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
