package service

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type UserService interface {
	CreateUser(input CreateUserInput) error
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *userService {
	return &userService{userRepo: userRepo}
}

type UserRepository interface {
	NextID() uuid.UUID
	Store(user model.User) error
	List(userIDs []model.UserID) ([]model.User, error)
}

type CreateUserInput struct {
	Login    string
	Password string
	AboutMe  string
}

func (service *userService) CreateUser(input CreateUserInput) error {
	user := model.NewUser(
		model.UserID(service.userRepo.NextID()),
		maybe.Nothing[model.ImageID](),
		input.Login,
		model.Client,
		input.Password,
		input.AboutMe,
	)

	err := service.userRepo.Store(user)
	if err != nil {
		return err
	}

	return nil
}
