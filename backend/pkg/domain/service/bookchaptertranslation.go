package service

import "server/pkg/domain/model"

type BookChapterTranslationService interface {
	StoreBookChapterTranslation(input StoreBookChapterTranslationInput) error
}

type bookChapterTranslationService struct {
	bookChapterTranslationRepo BookChapterTranslationRepository
}

func NewBookChapterTranslationService(
	bookChapterTranslationRepo BookChapterTranslationRepository,
) *bookChapterTranslationService {
	return &bookChapterTranslationService{
		bookChapterTranslationRepo: bookChapterTranslationRepo,
	}
}

type BookChapterTranslationRepository interface {
	Store(bookChapterTranslation model.BookChapterTranslation) error
}

type StoreBookChapterTranslationInput struct {
	BookChapterID model.BookChapterID
	TranslatorID  model.UserID
	Text          string
}

func (service *bookChapterTranslationService) StoreBookChapterTranslation(input StoreBookChapterTranslationInput) error {
	bookChapterTranslation := model.NewBookChapterTranslation(
		input.BookChapterID,
		input.TranslatorID,
		input.Text,
	)

	return service.bookChapterTranslationRepo.Store(bookChapterTranslation)
}
