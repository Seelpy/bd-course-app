package service

import (
	"server/pkg/domain/model"
)

type BookAuthorService interface {
	StoreBookAuthor(bookID model.BookID, authorID model.AuthorID) error
	DeleteBookAuthor(bookID model.BookID, authorID model.AuthorID) error
}

type bookAuthorService struct {
	bookAuthorRepo BookAuthorRepository
}

func NewBookAuthorService(bookAuthorRepo BookAuthorRepository) *bookAuthorService {
	return &bookAuthorService{bookAuthorRepo: bookAuthorRepo}
}

type BookAuthorRepository interface {
	Store(bookAuthor model.BookAuthor) error
	Delete(bookID model.BookID, authorID model.AuthorID) error
}

func (service *bookAuthorService) StoreBookAuthor(bookID model.BookID, authorID model.AuthorID) error {
	bookAuthor := model.NewBookAuthor(bookID, authorID)

	return service.bookAuthorRepo.Store(bookAuthor)
}

func (service *bookAuthorService) DeleteBookAuthor(bookID model.BookID, authorID model.AuthorID) error {
	return service.bookAuthorRepo.Delete(bookID, authorID)
}
