package service

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type BookService interface {
	CreateBook(input CreateBookInput) error
	EditBook(input EditBookInput) error
	EditBookImage(input EditBookImageInput) error
	PublishBook(input PublishBookInput) error
	DeleteBook(bookID model.BookID) error
}

type bookService struct {
	bookRepo BookRepository
}

func NewBookService(bookRepo BookRepository) *bookService {
	return &bookService{bookRepo: bookRepo}
}

type BookRepository interface {
	NextID() uuid.UUID
	Store(book model.Book) error
	Delete(bookID model.BookID) error
	FindByID(bookID model.BookID) (model.Book, error)
}

type CreateBookInput struct {
	Title       string
	Description string
}

type EditBookInput struct {
	ID          model.BookID
	Title       string
	Description string
}

type PublishBookInput struct {
	ID          model.BookID
	IsPublished bool
}

type EditBookImageInput struct {
	ID      model.BookID
	ImageID model.ImageID
}

func (service *bookService) CreateBook(input CreateBookInput) error {
	book := model.NewBook(
		model.BookID(service.bookRepo.NextID()),
		maybe.Nothing[model.ImageID](),
		input.Title,
		input.Description,
		false,
	)

	err := service.bookRepo.Store(book)
	if err != nil {
		return err
	}

	return nil
}

func (service *bookService) EditBook(input EditBookInput) error {
	book, err := service.bookRepo.FindByID(input.ID)
	if err != nil {
		return err
	}

	book.SetTitle(input.Title)
	book.SetDescription(input.Description)

	return service.bookRepo.Store(book)
}

func (service *bookService) EditBookImage(input EditBookImageInput) error {
	book, err := service.bookRepo.FindByID(input.ID)
	if err != nil {
		return err
	}

	book.SetCoverID(maybe.Just(input.ImageID))

	return service.bookRepo.Store(book)
}

func (service *bookService) PublishBook(input PublishBookInput) error {
	book, err := service.bookRepo.FindByID(input.ID)
	if err != nil {
		return err
	}

	book.SetIsPublished(input.IsPublished)

	return service.bookRepo.Store(book)
}

func (service *bookService) DeleteBook(bookID model.BookID) error {
	return service.bookRepo.Delete(bookID)
}
