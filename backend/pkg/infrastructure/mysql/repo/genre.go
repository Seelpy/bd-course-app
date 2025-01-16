package repo

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type GenreRepository struct {
	connection *sqlx.DB
}

func NewGenreRepository(connection *sqlx.DB) *GenreRepository {
	return &GenreRepository{connection: connection}
}

func (repo *GenreRepository) NextID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (repo *GenreRepository) Store(genre model.Genre) error {
	const query = `
		INSERT INTO
			genre (
				genre_id,
				name
			)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name)
	`

	binaryGenreID, err := uuid.UUID(genre.ID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query,
		binaryGenreID,
		genre.Name(),
	)

	return err
}

func (repo *GenreRepository) Delete(genreID model.GenreID) error {
	const query = `DELETE FROM genre WHERE genre_id = ?`

	binaryGenreID, err := uuid.UUID(genreID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryGenreID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		return model.ErrGenreNotFound
	}

	return err
}

func (repo *GenreRepository) FindByID(genreID model.GenreID) (model.Genre, error) {
	const query = `
		SELECT
			name
		FROM genre
		WHERE genre_id = ?
	`

	var name string
	binaryGenreID, err := uuid.UUID(genreID).MarshalBinary()
	if err != nil {
		return model.Genre{}, err
	}

	err = repo.connection.Get(&name, query, binaryGenreID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Genre{}, model.ErrGenreNotFound
	}
	if err != nil {
		return model.Genre{}, err
	}

	return model.NewGenre(
		genreID,
		name,
	), nil
}
