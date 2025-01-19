package query

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type AuthorQueryService interface {
	FindByID(authorID model.AuthorID) (AuthorOutput, error)
	List() ([]AuthorOutput, error)
	ListByBookID(bookID model.BookID) ([]AuthorOutput, error)
}

type AuthorOutput struct {
	AuthorID   uuid.UUID
	Avatar     maybe.Maybe[string]
	FirstName  string
	SecondName string
	MiddleName maybe.Maybe[string]
	Nickname   maybe.Maybe[string]
}

type authorQueryService struct {
	connection *sqlx.DB
}

func NewAuthorQueryService(connection *sqlx.DB) *authorQueryService {
	return &authorQueryService{connection: connection}
}

func (service *authorQueryService) FindByID(authorID model.AuthorID) (AuthorOutput, error) {
	const query = `
		SELECT 
			a.author_id,
			i.path AS avatar,
			a.first_name,
			a.second_name,
			a.middle_name,
			a.nickname
		FROM author a
		LEFT JOIN image i ON a.avatar_id = i.image_id
		WHERE a.author_id = ?;
	`

	binaryAuthorID, err := uuid.UUID(authorID).MarshalBinary()
	if err != nil {
		return AuthorOutput{}, err
	}

	var author sqlxAuthor
	err = service.connection.Get(&author, query, binaryAuthorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AuthorOutput{}, model.ErrAuthorNotFound
		}
		return AuthorOutput{}, err
	}

	avatar := maybe.Nothing[string]()
	if author.Avatar.Valid {
		avatar = maybe.Just(author.Avatar.String)
	}

	middleName := maybe.Nothing[string]()
	if author.MiddleName.Valid {
		middleName = maybe.Just(author.MiddleName.String)
	}

	nickname := maybe.Nothing[string]()
	if author.Nickname.Valid {
		nickname = maybe.Just(author.Nickname.String)
	}

	return AuthorOutput{
		AuthorID:   author.AuthorID,
		Avatar:     avatar,
		FirstName:  author.FirstName,
		SecondName: author.SecondName,
		MiddleName: middleName,
		Nickname:   nickname,
	}, nil
}

func (service *authorQueryService) List() ([]AuthorOutput, error) {
	const query = `
		SELECT 
			a.author_id,
			i.path AS avatar,
			a.first_name,
			a.second_name,
			a.middle_name,
			a.nickname
		FROM author a
		LEFT JOIN image i ON a.avatar_id = i.image_id;
	`

	var sqlxAuthors []sqlxAuthor
	err := service.connection.Select(&sqlxAuthors, query)
	if err != nil {
		return nil, err
	}

	authorOutputs := make([]AuthorOutput, len(sqlxAuthors))
	for i, a := range sqlxAuthors {
		avatar := maybe.Nothing[string]()
		if a.Avatar.Valid {
			avatar = maybe.Just(a.Avatar.String)
		}

		middleName := maybe.Nothing[string]()
		if a.MiddleName.Valid {
			middleName = maybe.Just(a.MiddleName.String)
		}

		nickname := maybe.Nothing[string]()
		if a.Nickname.Valid {
			nickname = maybe.Just(a.Nickname.String)
		}

		authorOutputs[i] = AuthorOutput{
			AuthorID:   a.AuthorID,
			Avatar:     avatar,
			FirstName:  a.FirstName,
			SecondName: a.SecondName,
			MiddleName: middleName,
			Nickname:   nickname,
		}
	}

	return authorOutputs, nil
}

func (service *authorQueryService) ListByBookID(bookID model.BookID) ([]AuthorOutput, error) {
	const query = `
		SELECT 
			a.author_id,
			i.path AS avatar,
			a.first_name,
			a.second_name,
			a.middle_name,
			a.nickname
		FROM author a
		INNER JOIN book_author ba ON a.author_id = ba.author_id AND ba.book_id = ?
		LEFT JOIN image i ON a.avatar_id = i.image_id;
	`

	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return nil, err
	}

	var sqlxAuthors []sqlxAuthor
	err = service.connection.Select(&sqlxAuthors, query, binaryBookID)
	if err != nil {
		return nil, err
	}

	authorOutputs := make([]AuthorOutput, len(sqlxAuthors))
	for i, a := range sqlxAuthors {
		avatar := maybe.Nothing[string]()
		if a.Avatar.Valid {
			avatar = maybe.Just(a.Avatar.String)
		}

		middleName := maybe.Nothing[string]()
		if a.MiddleName.Valid {
			middleName = maybe.Just(a.MiddleName.String)
		}

		nickname := maybe.Nothing[string]()
		if a.Nickname.Valid {
			nickname = maybe.Just(a.Nickname.String)
		}

		authorOutputs[i] = AuthorOutput{
			AuthorID:   a.AuthorID,
			Avatar:     avatar,
			FirstName:  a.FirstName,
			SecondName: a.SecondName,
			MiddleName: middleName,
			Nickname:   nickname,
		}
	}

	return authorOutputs, nil
}

type sqlxAuthor struct {
	AuthorID   uuid.UUID      `db:"author_id"`
	Avatar     sql.NullString `db:"avatar"`
	FirstName  string         `db:"first_name"`
	SecondName string         `db:"second_name"`
	MiddleName sql.NullString `db:"middle_name"`
	Nickname   sql.NullString `db:"nickname"`
}
