package repo

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
	"server/pkg/domain/service"
)

type bookRatingRepository struct {
	connection *sqlx.DB
}

func NewBookRatingRepository(connection *sqlx.DB) service.BookRatingRepository {
	return &bookRatingRepository{connection}
}

func (repo *bookRatingRepository) NextID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (repo *bookRatingRepository) Store(bookRating model.BookRating) error {
	const query = `
		INSERT INTO book_rating (
			book_id,
			user_id,
			value
		) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			value = VALUES(value)
	`

	binaryBookID, err := bookRating.BookID().MarshalBinary()
	if err != nil {
		return err
	}

	binaryUserID, err := uuid.UUID(bookRating.UserID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query, binaryBookID, binaryUserID, bookRating.Value)
	return err
}

func (repo *bookRatingRepository) Delete(bookID model.BookID, userID model.UserID) error {
	const query = `DELETE FROM book_rating WHERE book_id = ? AND user_id = ?`

	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return err
	}

	binaryUserID, err := uuid.UUID(userID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryBookID, binaryUserID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		return nil
	}

	return err
}

func (repo *bookRatingRepository) Find(bookID model.BookID, userID model.UserID) (model.BookRating, error) {
	const query = `
		SELECT value FROM book_rating WHERE book_id = ? AND user_id = ?
	`

	var ratingValue int
	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return model.BookRating{}, err
	}

	binaryUserID, err := uuid.UUID(userID).MarshalBinary()
	if err != nil {
		return model.BookRating{}, err
	}

	err = repo.connection.Get(&ratingValue, query, binaryBookID, binaryUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.BookRating{}, model.ErrBookRatingNotFound
	}
	if err != nil {
		return model.BookRating{}, err
	}

	return model.NewBookRating(bookID, userID, ratingValue), nil
}

func (repo *bookRatingRepository) AverageByBookID(bookID model.BookID) (float64, error) {
	const query = `
		SELECT AVG(value) FROM book_rating WHERE book_id = ?
	`

	var avg float64
	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return 0, err
	}

	err = repo.connection.Get(&avg, query, binaryBookID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return avg, nil
}

func (repo *bookRatingRepository) CountByBookID(bookID model.BookID) (int, error) {
	const query = `
		SELECT COUNT(*) FROM book_rating WHERE book_id = ?
	`

	var count int
	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return 0, err
	}

	err = repo.connection.Get(&count, query, binaryBookID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}
