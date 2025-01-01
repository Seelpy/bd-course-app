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
			 'role' = VALUES('role'),
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

	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) Delete(userID model.UserID) error {
	const query = `DELETE FROM user WHERE user_id = ?`

	result, err := repo.connection.Exec(query, userID)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) FindByID(userID model.UserID) (model.User, error) {
	const query = `
		SELECT
			user_id,
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
		model.UserID(user.ID),
		maybe.Nothing[model.ImageID](),
		user.Login,
		model.UserRole(user.Role),
		user.Password,
		user.AboutMe,
	), nil
}

type sqlxUser struct {
	ID       uuid.UUID `db:"user_id"`
	Login    string    `db:"login"`
	Role     int       `db:"role"`
	Password string    `db:"password"`
	AboutMe  string    `db:"about_me"`
}
