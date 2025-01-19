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
	List() ([]model.User, error)
	FindByLogin(login string) (model.User, error)
	FindByID(userID model2.UserID) (model.User, error)
}

type userQueryService struct {
	connection *sqlx.DB
}

func NewUserQueryService(connection *sqlx.DB) *userQueryService {
	return &userQueryService{connection}
}

func (service *userQueryService) List() ([]model.User, error) {
	const query = `
		SELECT
			u.user_id,
			i.path,
			u.login,
			u.role,
			u.password,
       		u.about_me
		FROM user u
		LEFT JOIN image i ON u.avatar_id = i.image_id;
`

	var sqlxUsers []sqlxUser
	err := service.connection.Select(&sqlxUsers, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model2.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	users := make([]model.User, len(sqlxUsers))
	for i, user := range sqlxUsers {
		var avatar maybe.Maybe[string]
		if user.AvatarID.Valid {
			avatar = maybe.Just(user.AvatarID.String)
		}

		users[i] = model.User{
			ID:       user.ID,
			Avatar:   avatar,
			Login:    user.Login,
			Role:     user.Role,
			Password: user.Password,
			AboutMe:  user.AboutMe,
		}
	}

	return users, nil
}

func (service *userQueryService) FindByLogin(login string) (model.User, error) {
	const query = `
		SELECT
			u.user_id,
			i.path,
			u.login,
			u.role,
			u.password,
       		u.about_me
		FROM user u
		LEFT JOIN image i ON u.avatar_id = i.image_id
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

	var avatar maybe.Maybe[string]
	if user.AvatarID.Valid {
		avatar = maybe.Just(user.AvatarID.String)
	}

	return model.User{
		ID:       user.ID,
		Avatar:   avatar,
		Login:    user.Login,
		Role:     user.Role,
		Password: user.Password,
		AboutMe:  user.AboutMe,
	}, nil
}

func (service *userQueryService) FindByID(userID model2.UserID) (model.User, error) {
	const query = `
		SELECT
			u.user_id,
			i.path,
			u.login,
			u.role,
			u.password,
       		u.about_me
		FROM user u
		LEFT JOIN image i ON u.avatar_id = i.image_id
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

	var avatar maybe.Maybe[string]
	if user.AvatarID.Valid {
		avatar = maybe.Just(user.AvatarID.String)
	}

	return model.User{
		ID:       user.ID,
		Avatar:   avatar,
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
