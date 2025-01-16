package query

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
	"time"
)

type VerifyBookRequestQueryService interface {
	List() ([]VerifyBookRequestOutput, error)
}

type VerifyBookRequestOutput struct {
	VerifyBookRequestID uuid.UUID
	TranslatorID        uuid.UUID
	BookID              uuid.UUID
	IsVerified          maybe.Maybe[bool]
	SendDate            time.Time
}

type verifyBookRequestQueryService struct {
	connection *sqlx.DB
}

func NewVerifyBookRequestQueryService(connection *sqlx.DB) *verifyBookRequestQueryService {
	return &verifyBookRequestQueryService{connection: connection}
}

func (service *verifyBookRequestQueryService) List() ([]VerifyBookRequestOutput, error) {
	const query = `
		SELECT *
		FROM verify_book_request vbr
		ORDER BY vbr.send_date DESC;
	`

	var sqlxVerifyBooksRequest []sqlxVerifyBookRequest
	err := service.connection.Select(&sqlxVerifyBooksRequest, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrBookChapterTranslationNotFound
	}
	if err != nil {
		return nil, err
	}

	verifyBookRequestsOutput := make([]VerifyBookRequestOutput, len(sqlxVerifyBooksRequest))
	for i, v := range sqlxVerifyBooksRequest {
		isVerified := maybe.Nothing[bool]()
		ok := v.IsVerified.Valid
		if ok {
			isVerified = maybe.Just(v.IsVerified.Bool)
		}

		verifyBookRequestsOutput[i] = VerifyBookRequestOutput{
			VerifyBookRequestID: v.VerifyBookRequestID,
			TranslatorID:        v.TranslatorID,
			BookID:              v.BookID,
			IsVerified:          isVerified,
			SendDate:            v.SendDate,
		}
	}

	return verifyBookRequestsOutput, nil
}

type sqlxVerifyBookRequest struct {
	VerifyBookRequestID uuid.UUID    `db:"verify_book_request_id"`
	TranslatorID        uuid.UUID    `db:"translator_id"`
	BookID              uuid.UUID    `db:"book_id"`
	IsVerified          sql.NullBool `db:"is_verified"`
	SendDate            time.Time    `db:"send_date"`
}
