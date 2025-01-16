package repo

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type BookChapterTranslationRepository struct {
	connection *sqlx.DB
}

func NewBookChapterTranslationRepository(connection *sqlx.DB) *BookChapterTranslationRepository {
	return &BookChapterTranslationRepository{connection}
}

func (repo *BookChapterTranslationRepository) Store(bookChapterTranslation model.BookChapterTranslation) error {
	const query = `
		INSERT INTO
			book_chapter_translation (
			      book_chapter_id,
			      translator_id,
			      text
			)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			text = VALUES(text)
	`

	binaryBookChapterID, err := uuid.UUID(bookChapterTranslation.BookChapterID()).MarshalBinary()
	if err != nil {
		return err
	}
	binaryTranslatorID, err := uuid.UUID(bookChapterTranslation.TranslatorID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query,
		binaryBookChapterID,
		binaryTranslatorID,
		bookChapterTranslation.Text(),
	)

	return err
}
