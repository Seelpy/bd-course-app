package service

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type UserService interface {
	CreateUser(input CreateUserInput) error
	EditUser(input EditUserInput) error
	EditImageUser(input EditUserImageInput) error
	DeleteUser(userID model.UserID) error
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
	Delete(userID model.UserID) error
	FindByID(userID model.UserID) (model.User, error)
}

type CreateUserInput struct {
	Login    string
	Password string
	AboutMe  string
}

type EditUserInput struct {
	ID       model.UserID
	Login    string
	Password string
	AboutMe  string
}

type EditUserImageInput struct {
	ID      model.UserID
	ImageID model.ImageID
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

	return service.userRepo.Store(user)
}

func (service *userService) EditUser(input EditUserInput) error {
	user, err := service.userRepo.FindByID(input.ID)
	if err != nil {
		return err
	}

	user.SetLogin(input.Login)
	user.SetPassword(input.Password)
	user.SetAboutMe(input.AboutMe)

	return service.userRepo.Store(user)
}

func (service *userService) EditImageUser(input EditUserImageInput) error {
	user, err := service.userRepo.FindByID(input.ID)
	if err != nil {
		return err
	}

	user.SetAvatarID(maybe.Just(input.ImageID))

	return service.userRepo.Store(user)
}

func (service *userService) DeleteUser(userID model.UserID) error {
	user, err := service.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user.Role() == model.Admin {
		return model.ErrNotDeleteAdmin
	}

	return service.userRepo.Delete(userID)
}
