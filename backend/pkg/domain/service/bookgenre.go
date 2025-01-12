package service

import (
	"server/pkg/domain/model"
)

type BookGenreService interface {
	StoreBookGenre(bookID model.BookID, genreID model.GenreID) error
	DeleteBookGenre(bookID model.BookID, genreID model.GenreID) error
}

type bookGenreService struct {
	bookGenreRepo BookGenreRepository
}

func NewBookGenreService(bookGenreRepo BookGenreRepository) *bookGenreService {
	return &bookGenreService{bookGenreRepo: bookGenreRepo}
}

type BookGenreRepository interface {
	Store(book model.BookGenre) error
	Delete(bookID model.BookID, genreID model.GenreID) error
}

func (service *bookGenreService) StoreBookGenre(bookID model.BookID, genreID model.GenreID) error {
	bookGenre := model.NewBookGenre(bookID, genreID)

	return service.bookGenreRepo.Store(bookGenre)
}

func (service *bookGenreService) DeleteBookGenre(bookID model.BookID, genreID model.GenreID) error {
	return service.bookGenreRepo.Delete(bookID, genreID)
}
