package repo

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
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
		VALUES (
			?,
		    ?,
		    ?,
		    ?,
		    ?
		)
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

func (repo *UserRepository) List(userIDs []model.UserID) ([]model.User, error) {
	return nil, nil
}
