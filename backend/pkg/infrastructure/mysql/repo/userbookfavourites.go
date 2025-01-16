package repo

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type userBookFavouritesRepository struct {
	connection *sqlx.DB
}

func NewUserBookFavouritesRepository(connection *sqlx.DB) *userBookFavouritesRepository {
	return &userBookFavouritesRepository{connection: connection}
}

func (repo *userBookFavouritesRepository) Store(userBookFavourites model.UserBookFavourites) error {
	const query = `
		INSERT INTO
			user_book_favourites (
			      user_id,
			      book_id,
			      type
			)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			type = VALUES(type)
	`

	binaryUserID, err := uuid.UUID(userBookFavourites.UserID()).MarshalBinary()
	if err != nil {
		return err
	}
	binaryBookID, err := uuid.UUID(userBookFavourites.BookID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query,
		binaryUserID,
		binaryBookID,
		userBookFavourites.UserBookFavouritesType(),
	)

	return err
}

func (repo *userBookFavouritesRepository) Delete(userID model.UserID, bookID model.BookID) error {
	const query = `DELETE FROM user_book_favourites WHERE user_id = ? AND book_id = ?`

	binaryUserID, err := uuid.UUID(userID).MarshalBinary()
	if err != nil {
		return err
	}
	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryUserID, binaryBookID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		return model.ErrUserBookFavouritesNotFound
	}

	return err
}
