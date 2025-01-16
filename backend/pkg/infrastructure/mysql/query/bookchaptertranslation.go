package query

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type BookChapterTranslationQueryService interface {
	GetByBookChapterIDAndTranslatorID(bookChapterID model.BookChapterID, translatorID model.UserID) (BookChapterTranslationOutput, error)
	ListTranslatorsByBookChapterId(bookChapterID model.BookChapterID) ([]uuid.UUID, error)
}

type BookChapterTranslationOutput struct {
	Text string
}

type bookChapterTranslationQueryService struct {
	connection *sqlx.DB
}

func NewBookChapterTranslationQueryService(connection *sqlx.DB) *bookChapterTranslationQueryService {
	return &bookChapterTranslationQueryService{
		connection: connection,
	}
}

func (service *bookChapterTranslationQueryService) GetByBookChapterIDAndTranslatorID(
	bookChapterID model.BookChapterID,
	translatorID model.UserID,
) (BookChapterTranslationOutput, error) {
	const query = `
		SELECT bct.text
		FROM book_chapter_translation bct
		WHERE bct.book_chapter_id = ? AND bct.translator_id = ?
	`

	binaryBookChapterID, err := uuid.UUID(bookChapterID).MarshalBinary()
	if err != nil {
		return BookChapterTranslationOutput{}, err
	}

	binaryTranslatorID, err := uuid.UUID(translatorID).MarshalBinary()
	if err != nil {
		return BookChapterTranslationOutput{}, err
	}

	var text string
	err = service.connection.Get(&text, query, binaryBookChapterID, binaryTranslatorID)
	if errors.Is(err, sql.ErrNoRows) {
		return BookChapterTranslationOutput{}, model.ErrBookChapterTranslationNotFound
	}
	if err != nil {
		return BookChapterTranslationOutput{}, err
	}

	return BookChapterTranslationOutput{
		Text: text,
	}, nil
}

func (service *bookChapterTranslationQueryService) ListTranslatorsByBookChapterId(bookChapterID model.BookChapterID) ([]uuid.UUID, error) {
	const query = `
		SELECT bct.translator_id
		FROM book_chapter_translation bct
		WHERE bct.book_chapter_id = ?
	`

	binaryBookChapterID, err := uuid.UUID(bookChapterID).MarshalBinary()
	if err != nil {
		return nil, err
	}

	var sqlxBookChapterTranslators []sqlxBookChapterTranslator
	err = service.connection.Select(&sqlxBookChapterTranslators, query, binaryBookChapterID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrBookChapterTranslationNotFound
	}
	if err != nil {
		return nil, err
	}

	translatorsID := make([]uuid.UUID, len(sqlxBookChapterTranslators))
	for i, translator := range sqlxBookChapterTranslators {
		translatorsID[i] = translator.BookChapterID
	}

	return translatorsID, nil
}

type sqlxBookChapterTranslator struct {
	BookChapterID uuid.UUID `db:"translator_id"`
}
