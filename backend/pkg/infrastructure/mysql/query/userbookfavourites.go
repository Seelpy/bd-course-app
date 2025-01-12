package query

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
	"strings"
)

type UserBookFavouritesQueryService interface {
	ListUserBookFavouritesByBook(userID model.UserID, bookID model.BookID) (model.UserBookFavouritesType, error)
	ListBookByUserBookFavourites(userID model.UserID, userBookFavouritesTypes []model.UserBookFavouritesType) ([]UserBookFavouritesBooksOutput, error)
}

type userBookFavouritesQueryService struct {
	connection *sqlx.DB
}

func NewUserBookFavouritesQueryService(connection *sqlx.DB) *userBookFavouritesQueryService {
	return &userBookFavouritesQueryService{connection: connection}
}

type UserBookFavouritesBooksOutput struct {
	Type  model.UserBookFavouritesType
	Books []BookOutput
}

func (service *userBookFavouritesQueryService) ListUserBookFavouritesByBook(
	userID model.UserID,
	bookID model.BookID,
) (model.UserBookFavouritesType, error) {
	const query = `
		SELECT
			type
		FROM user_book_favourites
		WHERE user_id = ? AND book_id = ?;
`

	binaryUserID, err := uuid.UUID(userID).MarshalBinary()
	if err != nil {
		return model.READING, err
	}
	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return model.READING, err
	}

	var typeInt int
	err = service.connection.Get(&typeInt, query, binaryUserID, binaryBookID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.READING, model.ErrUserBookFavouritesNotFound
	}
	if err != nil {
		return model.READING, err
	}

	return model.UserBookFavouritesType(typeInt), nil
}

func (service *userBookFavouritesQueryService) ListBookByUserBookFavourites(
	userID model.UserID,
	userBookFavouritesTypes []model.UserBookFavouritesType,
) ([]UserBookFavouritesBooksOutput, error) {
	if len(userBookFavouritesTypes) == 0 {
		return nil, nil
	}

	query := `
		SELECT
			ubf.type,
			b.book_id,
			i.path,
			b.title,
			b.description
		FROM user_book_favourites ubf
		LEFT JOIN book b ON ubf.book_id = b.id
		LEFT JOIN image i ON b.cover_id = i.image_id
		WHERE user_id = ?;
`
	binaryUserID, err := uuid.UUID(userID).MarshalBinary()
	if err != nil {
		return nil, err
	}

	query += " AND IN (?)"
	values := make([]string, len(userBookFavouritesTypes))
	for i, v := range userBookFavouritesTypes {
		values[i] = fmt.Sprintf("%v", v)
	}
	valuesStr := strings.Join(values, ", ")
	query = strings.Replace(query, "?", valuesStr, 1)

	var userBookFavouritesBooksOutput []sqlxUserBookFavouritesBooksOutput
	err = service.connection.Select(&userBookFavouritesBooksOutput, query, binaryUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrUserBookFavouritesNotFound
	}
	if err != nil {
		return nil, err
	}

	outputs := convertToUserBookFavouritesBooksOutput(userBookFavouritesBooksOutput)
	return outputs, nil
}

func convertToUserBookFavouritesBooksOutput(sqlxOutputs []sqlxUserBookFavouritesBooksOutput) []UserBookFavouritesBooksOutput {
	grouped := make(map[model.UserBookFavouritesType][]BookOutput)
	for _, output := range sqlxOutputs {
		cover := maybe.Nothing[string]()
		if output.Cover.Valid {
			cover = maybe.Just(output.Cover.String)
		}

		book := BookOutput{
			BookID:      output.BookID,
			Cover:       cover,
			Title:       output.Title,
			Description: output.Description,
		}
		grouped[model.UserBookFavouritesType(output.Type)] = append(
			grouped[model.UserBookFavouritesType(output.Type)],
			book,
		)
	}

	var result []UserBookFavouritesBooksOutput
	for t, books := range grouped {
		result = append(result, UserBookFavouritesBooksOutput{
			Type:  t,
			Books: books,
		})
	}

	return result
}

type sqlxUserBookFavouritesBooksOutput struct {
	Type        int            `db:"type"`
	BookID      uuid.UUID      `db:"book_id"`
	Cover       sql.NullString `db:"path"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
}
