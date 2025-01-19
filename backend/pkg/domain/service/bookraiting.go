package service

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type BookRatingService interface {
	StoreRating(input StoreBookRatingInput) error
	DeleteRating(bookID model.BookID, userID model.UserID) error
	GetStatistics(bookID model.BookID, userID maybe.Maybe[model.UserID]) (StatisticsBookRatingOutput, error)
}

type bookRatingService struct {
	bookRatingRepo BookRatingRepository
}

func NewBookRatingService(bookRatingRepo BookRatingRepository) BookRatingService {
	return &bookRatingService{bookRatingRepo: bookRatingRepo}
}

type BookRatingRepository interface {
	NextID() uuid.UUID
	Store(bookRating model.BookRating) error
	Delete(bookID model.BookID, userID model.UserID) error
	Find(bookID model.BookID, userID model.UserID) (model.BookRating, error)
	AverageByBookID(bookID model.BookID) (float64, error)
	CountByBookID(bookID model.BookID) (int, error)
}

type StoreBookRatingInput struct {
	BookID model.BookID
	UserID model.UserID
	Value  int
}

type StatisticsBookRatingOutput struct {
	Average         float64
	Count           int
	UserLoginRating maybe.Maybe[int]
}

func (service *bookRatingService) StoreRating(input StoreBookRatingInput) error {
	bookRating := model.NewBookRating(
		input.BookID,
		input.UserID,
		input.Value,
	)

	return service.bookRatingRepo.Store(bookRating)
}

func (service *bookRatingService) DeleteRating(bookID model.BookID, userID model.UserID) error {
	return service.bookRatingRepo.Delete(bookID, userID)
}

func (service *bookRatingService) GetStatistics(
	bookID model.BookID,
	userID maybe.Maybe[model.UserID],
) (StatisticsBookRatingOutput, error) {
	average, err := service.bookRatingRepo.AverageByBookID(bookID)
	if err != nil {
		return StatisticsBookRatingOutput{}, err
	}

	count, err := service.bookRatingRepo.CountByBookID(bookID)
	if err != nil {
		return StatisticsBookRatingOutput{}, err
	}

	userLoginRating := maybe.Nothing[int]()
	if userID2, ok := userID.Get(); ok {
		bookRating, err2 := service.bookRatingRepo.Find(bookID, userID2)
		if err2 != nil {
			return StatisticsBookRatingOutput{}, err2
		}

		userLoginRating = maybe.Just(bookRating.Value())
	}

	return StatisticsBookRatingOutput{
		Average:         average,
		Count:           count,
		UserLoginRating: userLoginRating,
	}, nil
}
