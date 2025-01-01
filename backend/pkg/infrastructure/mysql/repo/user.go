package repo

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type UserRepository struct {
	connection *sqlx.DB
}

func NewUserRepository(connection *sqlx.DB) *UserRepository {
	return &UserRepository{connection}
}

func (repo *UserRepository) NextID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (repo *UserRepository) Store(user model.User) error {
	const query = `
		INSERT INTO
			user (
			      user_id,
			      login,
			      role,
			      password,
			      about_me
			)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			 login = VALUES(login),
			 role = VALUES(role),
			 password = VALUES(password),
			 about_me = VALUES(about_me)
	`

	binaryUserID, err := uuid.UUID(user.ID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query,
		binaryUserID,
		user.Login(),
		user.Role(),
		user.Password(),
		user.AboutMe(),
	)

	return err
}

func (repo *UserRepository) Delete(userID model.UserID) error {
	const query = `DELETE FROM user WHERE user_id = ?`

	binaryUserID, err := uuid.UUID(userID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryUserID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		return model.ErrUserNotFound
	}

	return err
}

func (repo *UserRepository) FindByID(userID model.UserID) (model.User, error) {
	const query = `
		SELECT
			login,
			role,
			password,
			about_me
		FROM user
		WHERE user_id = ?
`

	var user sqlxUser
	binaryUserID, err := uuid.UUID(userID).MarshalBinary()
	if err != nil {
		return model.User{}, err
	}

	err = repo.connection.Get(&user, query, binaryUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, model.ErrUserNotFound
	}
	if err != nil {
		return model.User{}, err
	}

	return model.NewUser(
		model.UserID(userID),
		maybe.Nothing[model.ImageID](),
		user.Login,
		model.UserRole(user.Role),
		user.Password,
		user.AboutMe,
	), nil
}

type sqlxUser struct {
	Login    string `db:"login"`
	Role     int    `db:"role"`
	Password string `db:"password"`
	AboutMe  string `db:"about_me"`
}
