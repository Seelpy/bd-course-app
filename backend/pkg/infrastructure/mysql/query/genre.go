package query

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type GenreQueryService interface {
	List() ([]GenreOutput, error)
	ListByBookID(bookID model.BookID) ([]GenreOutput, error)
}

type GenreOutput struct {
	GenreID uuid.UUID `json:"id"`
	Name    string    `json:"name"`
}

type genreQueryService struct {
	connection *sqlx.DB
}

func NewGenreQueryService(connection *sqlx.DB) *genreQueryService {
	return &genreQueryService{connection: connection}
}

func (service *genreQueryService) List() ([]GenreOutput, error) {
	const query = `
		SELECT 
			genre_id,
			name
		FROM genre;
	`

	var sqlxGenres []sqlxGenre
	err := service.connection.Select(&sqlxGenres, query)
	if err != nil {
		return nil, err
	}

	genreOutputs := make([]GenreOutput, len(sqlxGenres))
	for i, g := range sqlxGenres {
		genreOutputs[i] = GenreOutput{
			GenreID: g.GenreID,
			Name:    g.Name,
		}
	}

	return genreOutputs, nil
}

func (service *genreQueryService) ListByBookID(bookID model.BookID) ([]GenreOutput, error) {
	const query = `
		SELECT 
			genre_id,
			name
		FROM genre g
		LEFT JOIN book_genre bg ON bg.genre_id = g.genre_id
		WHERE bg.book_id = ?;
	`

	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return nil, err
	}

	var sqlxGenres []sqlxGenre
	err = service.connection.Select(&sqlxGenres, query, binaryBookID)
	if err != nil {
		return nil, err
	}

	genreOutputs := make([]GenreOutput, len(sqlxGenres))
	for i, g := range sqlxGenres {
		genreOutputs[i] = GenreOutput{
			GenreID: g.GenreID,
			Name:    g.Name,
		}
	}

	return genreOutputs, nil
}

type sqlxGenre struct {
	GenreID uuid.UUID `db:"genre_id"`
	Name    string    `db:"name"`
}
