package repo

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type bookGenreRepository struct {
	connection *sqlx.DB
}

func NewBookGenreRepository(connection *sqlx.DB) *bookGenreRepository {
	return &bookGenreRepository{connection: connection}
}

func (repo *bookGenreRepository) Store(bookGenre model.BookGenre) error {
	const query = `
		INSERT INTO book_genre (
			book_id,
			genre_id
		)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			book_id = VALUES(book_id),
			genre_id = VALUES(genre_id)
	`

	binaryBookID, err := uuid.UUID(bookGenre.BookID()).MarshalBinary()
	if err != nil {
		return err
	}

	binaryGenreID, err := uuid.UUID(bookGenre.GenreID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query, binaryBookID, binaryGenreID)
	return err
}

func (repo *bookGenreRepository) Delete(bookID model.BookID, genreID model.GenreID) error {
	const query = `DELETE FROM book_genre WHERE book_id = ? AND genre_id = ?`

	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return err
	}

	binaryGenreID, err := uuid.UUID(genreID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryBookID, binaryGenreID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		return model.ErrBookGenreNotFound
	}

	return err
}
