package query

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mono83/maybe"
	"server/pkg/infrastructure/model"
)

type UserQueryService interface {
	FindByLogin(login string) (model.User, error)
}

type userQueryService struct {
	connection *sqlx.DB
}

func NewUserQueryService(connection *sqlx.DB) *userQueryService {
	return &userQueryService{connection}
}

func (service *userQueryService) FindByLogin(login string) (model.User, error) {
	const query = `
		SELECT
			user_id,
			avatar_id,
			login,
			role,
			password,
       		about_me
		FROM user
		WHERE login = ?
`

	var sqlxUser sqlxUser

	row := service.connection.QueryRow(query, login)

	err := row.Scan(
		&sqlxUser.ID,
		&sqlxUser.AvatarID,
		&sqlxUser.Login,
		&sqlxUser.Role,
		&sqlxUser.Password,
		&sqlxUser.AboutMe,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, nil
		}
		return model.User{}, err
	}

	var avatarID maybe.Maybe[uuid.UUID]
	if sqlxUser.AvatarID.Valid {
		id, err2 := uuid.FromString(sqlxUser.AvatarID.String)
		if err2 == nil {
			avatarID = maybe.Just(id)
		}
	}

	return model.User{
		ID:       sqlxUser.ID,
		AvatarID: avatarID,
		Login:    sqlxUser.Login,
		Role:     sqlxUser.Role,
		Password: sqlxUser.Password,
		AboutMe:  sqlxUser.AboutMe,
	}, nil
}

type sqlxUser struct {
	ID       uuid.UUID      `db:"user_id"`
	AvatarID sql.NullString `db:"avatar_id"`
	Login    string         `db:"login"`
	Role     int            `db:"role"`
	Password string         `db:"password"`
	AboutMe  string         `db:"about_me"`
}
