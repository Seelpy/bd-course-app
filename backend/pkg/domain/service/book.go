package service

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type BookService interface {
	CreateBook(input CreateBookInput) error
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
	List(bookIDs []model.BookID) ([]model.Book, error)
}

type CreateBookInput struct {
	Title       string
	Description string
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
