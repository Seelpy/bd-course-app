package provider

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type VerifyBookRequestProvider interface {
	FindBookIDByVerifyBookRequestID(verifyBookRequestID model.VerifyBookRequestID) (model.BookID, error)
}

type verifyBookRequestProvider struct {
	connection *sqlx.DB
}

func NewVerifyBookRequestProvider(connection *sqlx.DB) *verifyBookRequestProvider {
	return &verifyBookRequestProvider{connection}
}

func (provider *verifyBookRequestProvider) FindBookIDByVerifyBookRequestID(verifyBookRequestID model.VerifyBookRequestID) (model.BookID, error) {
	const query = `
		SELECT
			book_id
		FROM verify_book_request
		WHERE verify_book_request_id = ?
`

	var bookID uuid.UUID
	binaryVerifyBookRequestID, err := uuid.UUID(verifyBookRequestID).MarshalBinary()
	if err != nil {
		return model.BookID{}, err
	}

	err = provider.connection.Get(&bookID, query, binaryVerifyBookRequestID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.BookID{}, model.ErrVerifyBookRequestNotFound
	}
	if err != nil {
		return model.BookID{}, err
	}

	return model.BookID(bookID), nil
}
