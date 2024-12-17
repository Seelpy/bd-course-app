package repo

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
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
		VALUES (
			?,
		    ?,
		    ?,
		    ?
		)
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

func (repo *BookRepository) List(bookIDs []model.BookID) ([]model.Book, error) {
	return nil, nil
}
