package repo

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type BookChapterRepository struct {
	connection *sqlx.DB
}

func NewBookChapterRepository(connection *sqlx.DB) *BookChapterRepository {
	return &BookChapterRepository{connection}
}

func (repo *BookChapterRepository) NextID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (repo *BookChapterRepository) Store(bookChapter model.BookChapter) error {
	const query = `
		INSERT INTO
			book_chapter (
			      book_chapter_id,
			      book_id,
			      chapter_index,
			      title
			)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			title = VALUES(title),
			book_id = VALUES(book_id),
			chapter_index = VALUES(chapter_index)
	`

	binaryBookChapterID, err := uuid.UUID(bookChapter.ID()).MarshalBinary()
	if err != nil {
		return err
	}
	binaryBookID, err := uuid.UUID(bookChapter.BookID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query,
		binaryBookChapterID,
		binaryBookID,
		bookChapter.Index(),
		bookChapter.Title(),
	)

	return err
}

func (repo *BookChapterRepository) Delete(bookChapterID model.BookChapterID) error {
	const query = `DELETE FROM book_chapter WHERE book_chapter_id = ?`

	binaryBookChapterID, err := uuid.UUID(bookChapterID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryBookChapterID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		return model.ErrBookChapterNotFound
	}

	return err
}

func (repo *BookChapterRepository) FindByID(bookChapterID model.BookChapterID) (model.BookChapter, error) {
	const query = `
		SELECT
			book_id,
			chapter_index,
			title
		FROM book_chapter
		WHERE book_chapter_id = ?
`

	var bookChapter sqlxBookChapter
	binaryBookChapterID, err := uuid.UUID(bookChapterID).MarshalBinary()
	if err != nil {
		return model.BookChapter{}, err
	}

	err = repo.connection.Get(&bookChapter, query, binaryBookChapterID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.BookChapter{}, model.ErrBookChapterNotFound
	}
	if err != nil {
		return model.BookChapter{}, err
	}

	return model.NewBookChapter(
		bookChapterID,
		model.BookID(bookChapter.BookID),
		bookChapter.Index,
		bookChapter.Title,
	), nil
}

func (repo *BookChapterRepository) ListOrderIndexesByBookID(bookID model.BookID) ([]model.BookChapter, error) {
	const query = `
		SELECT
		    book_chapter_id,
			book_id,
			chapter_index,
			title
		FROM book_chapter
		WHERE book_id = ?
		ORDER BY chapter_index ASC
`

	var sqlxBookChapters []sqlxBookChapter
	binaryBookID, err := uuid.UUID(bookID).MarshalBinary()
	if err != nil {
		return nil, err
	}

	err = repo.connection.Select(&sqlxBookChapters, query, binaryBookID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrBookChapterNotFound
	}
	if err != nil {
		return nil, err
	}

	var bookChapters []model.BookChapter
	for _, bookChapter := range sqlxBookChapters {
		bookChapters = append(bookChapters, model.NewBookChapter(
			bookChapter.BookChapterID,
			bookChapter.BookID,
			bookChapter.Index,
			bookChapter.Title,
		))
	}

	return bookChapters, nil
}

type sqlxBookChapter struct {
	BookChapterID uuid.UUID `db:"book_chapter_id"`
	BookID        uuid.UUID `db:"book_id"`
	Index         int       `db:"chapter_index"`
	Title         string    `db:"title"`
}
