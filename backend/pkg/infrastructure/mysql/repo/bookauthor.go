package repo

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type bookAuthorRepository struct {
	connection *sqlx.DB
}

func NewBookAuthorRepository(connection *sqlx.DB) *bookAuthorRepository {
	return &bookAuthorRepository{connection: connection}
}

func (repo *bookAuthorRepository) Store(bookAuthor model.BookAuthor) error {
	const query = `
		INSERT INTO book_author (
			book_id,
			author_id
		)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			book_id = VALUES(book_id),
			author_id = VALUES(author_id)
	`

	binaryBookID, err := uuid.UUID(bookAuthor.BookID()).MarshalBinary()
	if err != nil {
		return err
	}

	binaryAuthorID, err := uuid.UUID(bookAuthor.Author()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query, binaryBookID, binaryAuthorID)
	return err
}

func (repo *bookAuthorRepository) Delete(bookID model.BookID, authorID model.AuthorID) error {
	const query = `DELETE FROM book_author WHERE book_id = ? AND author_id = ?`

	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return err
	}

	binaryAuthorID, err := uuid.UUID(authorID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryBookID, binaryAuthorID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		return model.ErrBookAuthorNotFound
	}

	return err
}
