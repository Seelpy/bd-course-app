package service

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
)

type AuthorService interface {
	CreateAuthor(input CreateAuthorInput) error
	EditAuthor(input EditAuthorInput) error
	EditAuthorAvatar(input EditAuthorAvatarInput) error
	DeleteAuthor(authorID model.AuthorID) error
}

type authorService struct {
	authorRepo AuthorRepository
}

func NewAuthorService(authorRepo AuthorRepository) *authorService {
	return &authorService{authorRepo: authorRepo}
}

type AuthorRepository interface {
	NextID() uuid.UUID
	Store(author model.Author) error
	Delete(authorID model.AuthorID) error
	FindByID(authorID model.AuthorID) (model.Author, error)
}

type CreateAuthorInput struct {
	FirstName  string
	SecondName string
	MiddleName maybe.Maybe[string]
	Nickname   maybe.Maybe[string]
}

type EditAuthorInput struct {
	ID         model.AuthorID
	FirstName  string
	SecondName string
	MiddleName maybe.Maybe[string]
	Nickname   maybe.Maybe[string]
}

type EditAuthorAvatarInput struct {
	ID      model.AuthorID
	ImageID model.ImageID
}

func (service *authorService) CreateAuthor(input CreateAuthorInput) error {
	author := model.NewAuthor(
		model.AuthorID(service.authorRepo.NextID()),
		maybe.Nothing[model.ImageID](),
		input.FirstName,
		input.SecondName,
		input.MiddleName,
		input.Nickname,
	)

	return service.authorRepo.Store(author)
}

func (service *authorService) EditAuthor(input EditAuthorInput) error {
	author, err := service.authorRepo.FindByID(input.ID)
	if err != nil {
		return err
	}

	author.SetFirstName(input.FirstName)
	author.SetSecondName(input.SecondName)
	author.SetMiddleName(input.MiddleName)
	author.SetNickname(input.Nickname)

	return service.authorRepo.Store(author)
}

func (service *authorService) EditAuthorAvatar(input EditAuthorAvatarInput) error {
	author, err := service.authorRepo.FindByID(input.ID)
	if err != nil {
		return err
	}

	author.SetAvatarID(maybe.Just(input.ImageID))

	return service.authorRepo.Store(author)
}

func (service *authorService) DeleteAuthor(authorID model.AuthorID) error {
	return service.authorRepo.Delete(authorID)
}
