package service

import (
	"github.com/gofrs/uuid"
	"server/pkg/domain/model"
)

type BookChapterService interface {
	CreateBookChapter(input CreateBookChapterInput) error
	EditBookChapter(input EditBookChapterInput) error
	DeleteBookChapter(bookChapterID uuid.UUID) error
}

type bookChapterService struct {
	bookChapterRepo BookChapterRepository
}

func NewBookChapterService(bookChapterRepo BookChapterRepository) *bookChapterService {
	return &bookChapterService{
		bookChapterRepo: bookChapterRepo,
	}
}

type BookChapterRepository interface {
	NextID() uuid.UUID
	Store(bookChapter model.BookChapter) error
	Delete(bookChapterID model.BookChapterID) error
	FindByID(bookChapterID model.BookChapterID) (model.BookChapter, error)
	ListOrderIndexesByBookID(bookID model.BookID) ([]model.BookChapter, error)
}

type CreateBookChapterInput struct {
	BookID model.BookID
	Title  string
}

type EditBookChapterInput struct {
	BookChapterID model.BookChapterID
	Title         string
}

func (service *bookChapterService) CreateBookChapter(input CreateBookChapterInput) error {
	indexes, err := service.bookChapterRepo.ListOrderIndexesByBookID(input.BookID)
	if err != nil {
		return err
	}
	nextOrder := 0
	if len(indexes) > 0 {
		nextOrder = indexes[len(indexes)-1].Index() + 1
	}

	bookChapter := model.NewBookChapter(
		model.BookChapterID(service.bookChapterRepo.NextID()),
		input.BookID,
		nextOrder,
		input.Title,
	)

	return service.bookChapterRepo.Store(bookChapter)
}

func (service *bookChapterService) EditBookChapter(input EditBookChapterInput) error {
	bookChapter, err := service.bookChapterRepo.FindByID(input.BookChapterID)
	if err != nil {
		return err
	}

	bookChapter.SetTitle(input.Title)

	return service.bookChapterRepo.Store(bookChapter)
}

func (service *bookChapterService) DeleteBookChapter(bookChapterID uuid.UUID) error {
	bookChapter, err := service.bookChapterRepo.FindByID(bookChapterID)
	if err != nil {
		return err
	}
	bookChapters, err := service.bookChapterRepo.ListOrderIndexesByBookID(bookChapter.BookID())
	if err != nil {
		return err
	}

	err = service.bookChapterRepo.Delete(bookChapterID)
	if err != nil {
		return err
	}

	var isDec bool
	for i, bookChap := range bookChapters {
		if bookChap.ID() == bookChapterID {
			isDec = true
			continue
		}

		newIndex := i
		if isDec {
			newIndex--
		}

		bookChap.SetIndex(newIndex)

		err = service.bookChapterRepo.Store(bookChap)
		if err != nil {
			return err
		}
	}

	return nil
}
