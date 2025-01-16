package service

import "server/pkg/domain/model"

type UserBookFavouritesService interface {
	StoreUserBookFavourites(input StoreUserBookFavouritesInput) error
	DeleteUserBookFavourites(input DeleteUserBookFavouritesInput) error
}

type userBookFavouritesService struct {
	userBookFavouritesRepository UserBookFavouritesRepository
}

func NewUserBookFavouritesService(userBookFavouritesRepository UserBookFavouritesRepository) *userBookFavouritesService {
	return &userBookFavouritesService{
		userBookFavouritesRepository: userBookFavouritesRepository,
	}
}

type UserBookFavouritesRepository interface {
	Store(userBookFavourites model.UserBookFavourites) error
	Delete(userID model.UserID, bookID model.BookID) error
}

type StoreUserBookFavouritesInput struct {
	UserID model.UserID
	BookID model.BookID
	Type   model.UserBookFavouritesType
}

type DeleteUserBookFavouritesInput struct {
	UserID model.UserID
	BookID model.BookID
}

func (service *userBookFavouritesService) StoreUserBookFavourites(input StoreUserBookFavouritesInput) error {
	userBookFavourites := model.NewUserBookFavourites(
		input.UserID,
		input.BookID,
		input.Type,
	)

	return service.userBookFavouritesRepository.Store(userBookFavourites)
}

func (service *userBookFavouritesService) DeleteUserBookFavourites(input DeleteUserBookFavouritesInput) error {
	return service.userBookFavouritesRepository.Delete(input.UserID, input.BookID)
}
