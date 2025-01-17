package query

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mono83/maybe"
	model2 "server/pkg/domain/model"
	"server/pkg/infrastructure/model"
)

type UserQueryService interface {
	FindByLogin(login string) (model.User, error)
	FindByID(userID model2.UserID) (model.User, error)
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
		WHERE login = ?;
`

	var user sqlxUser
	err := service.connection.Get(&user, query, login)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, model2.ErrUserNotFound
	}
	if err != nil {
		return model.User{}, err
	}

	var avatarID maybe.Maybe[uuid.UUID]
	if user.AvatarID.Valid {
		id, err2 := uuid.FromString(user.AvatarID.String)
		if err2 == nil {
			avatarID = maybe.Just(id)
		}
	}

	return model.User{
		ID:       user.ID,
		AvatarID: avatarID,
		Login:    user.Login,
		Role:     user.Role,
		Password: user.Password,
		AboutMe:  user.AboutMe,
	}, nil
}

func (service *userQueryService) FindByID(userID model2.UserID) (model.User, error) {
	const query = `
		SELECT
			user_id,
			avatar_id,
			login,
			role,
			password,
       		about_me
		FROM user
		WHERE user_id = ?;
`

	binaryUserID, err := uuid.UUID(userID).MarshalBinary()
	if err != nil {
		return model.User{}, err
	}

	var user sqlxUser
	err = service.connection.Get(&user, query, binaryUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, model2.ErrUserNotFound
	}
	if err != nil {
		return model.User{}, err
	}

	var avatarID maybe.Maybe[uuid.UUID]
	if user.AvatarID.Valid {
		id, err2 := uuid.FromString(user.AvatarID.String)
		if err2 == nil {
			avatarID = maybe.Just(id)
		}
	}

	return model.User{
		ID:       user.ID,
		AvatarID: avatarID,
		Login:    user.Login,
		Role:     user.Role,
		Password: user.Password,
		AboutMe:  user.AboutMe,
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
