package service

import (
	"github.com/gofrs/uuid"
	"server/pkg/domain/model"
)

type GenreService interface {
	CreateGenre(input CreateGenreInput) error
	EditGenre(input EditGenreInput) error
	DeleteGenre(genreID model.GenreID) error
}

type genreService struct {
	genreRepo GenreRepository
}

func NewGenreService(genreRepo GenreRepository) *genreService {
	return &genreService{genreRepo: genreRepo}
}

type GenreRepository interface {
	NextID() uuid.UUID
	Store(genre model.Genre) error
	Delete(genreID model.GenreID) error
	FindByID(genreID model.GenreID) (model.Genre, error)
}

type CreateGenreInput struct {
	Name string
}

type EditGenreInput struct {
	ID   model.GenreID
	Name string
}

func (service *genreService) CreateGenre(input CreateGenreInput) error {
	genre := model.NewGenre(
		model.GenreID(service.genreRepo.NextID()),
		input.Name,
	)

	return service.genreRepo.Store(genre)
}

func (service *genreService) EditGenre(input EditGenreInput) error {
	genre, err := service.genreRepo.FindByID(input.ID)
	if err != nil {
		return err
	}

	genre.SetName(input.Name)

	return service.genreRepo.Store(genre)
}

func (service *genreService) DeleteGenre(genreID model.GenreID) error {
	return service.genreRepo.Delete(genreID)
}
