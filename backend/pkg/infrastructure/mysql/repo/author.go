package repo

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type AuthorRepository struct {
	connection *sqlx.DB
}

func NewAuthorRepository(connection *sqlx.DB) *AuthorRepository {
	return &AuthorRepository{connection: connection}
}

func (repo *AuthorRepository) NextID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (repo *AuthorRepository) Store(author model.Author) error {
	const query = `
		INSERT INTO
			author (
				author_id,
				avatar_id,
				first_name,
				second_name,
				middle_name,
				nickname
			)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			avatar_id = VALUES(avatar_id),
			first_name = VALUES(first_name),
			second_name = VALUES(second_name),
			middle_name = VALUES(middle_name),
			nickname = VALUES(nickname)
	`

	binaryAuthorID, err := uuid.UUID(author.ID()).MarshalBinary()
	if err != nil {
		return err
	}

	var binaryAvatarID *[]byte
	if avatarID, ok := author.AvatarID().Get(); ok {
		uid, err2 := uuid.UUID(avatarID).MarshalBinary()
		if err2 != nil {
			return err2
		}
		binaryAvatarID = &uid
	} else {
		binaryAvatarID = nil
	}

	// Handle Maybe[string] for middle name
	var middleName *string
	if mn, ok := author.MiddleName().Get(); ok {
		middleName = &mn
	} else {
		middleName = nil
	}

	// Handle Maybe[string] for middle name
	var nickName *string
	if nn, ok := author.Nickname().Get(); ok {
		nickName = &nn
	} else {
		nickName = nil
	}

	_, err = repo.connection.Exec(query,
		binaryAuthorID,
		binaryAvatarID,
		author.FirstName(),
		author.SecondName(),
		middleName,
		nickName,
	)

	return err
}

func (repo *AuthorRepository) Delete(authorID model.AuthorID) error {
	const query = `DELETE FROM author WHERE author_id = ?`

	binaryAuthorID, err := uuid.UUID(authorID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryAuthorID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		return model.ErrAuthorNotFound
	}

	return err
}

func (repo *AuthorRepository) FindByID(authorID model.AuthorID) (model.Author, error) {
	const query = `
		SELECT
			first_name,
			second_name,
			middle_name,
			nickname
		FROM author
		WHERE author_id = ?
	`

	var author sqlxAuthor
	binaryAuthorID, err := uuid.UUID(authorID).MarshalBinary()
	if err != nil {
		return model.Author{}, err
	}

	err = repo.connection.Get(&author, query, binaryAuthorID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Author{}, model.ErrAuthorNotFound
	}
	if err != nil {
		return model.Author{}, err
	}

	middleName := maybe.Nothing[string]()
	if author.MiddleName.Valid {
		middleName = maybe.Just(author.MiddleName.String)
	}

	nickName := maybe.Nothing[string]()
	if author.Nickname.Valid {
		nickName = maybe.Just(author.Nickname.String)
	}

	return model.NewAuthor(
		authorID,
		maybe.Nothing[model.ImageID](),
		author.FirstName,
		author.SecondName,
		middleName,
		nickName,
	), nil
}

// sqlxAuthor представляет структуру автора, используемую для сканирования данных из базы данных
type sqlxAuthor struct {
	FirstName  string         `db:"first_name"`
	SecondName string         `db:"second_name"`
	MiddleName sql.NullString `db:"middle_name"`
	Nickname   sql.NullString `db:"nickname"`
}
