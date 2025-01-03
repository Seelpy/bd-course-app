package repo

import (
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
	"time"
)

type VerifyBookRequestRepository struct {
	connection *sqlx.DB
}

func NewVerifyBookRequestRepository(connection *sqlx.DB) *VerifyBookRequestRepository {
	return &VerifyBookRequestRepository{connection: connection}
}

func (repo *VerifyBookRequestRepository) NextID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (repo *VerifyBookRequestRepository) Store(verifyBookRequest model.VerifyBookRequest) error {
	const query = `
		INSERT INTO
			verify_book_request (
			      verify_book_request_id,
			      translator_id,
			      book_id,
			      is_verified,
			      send_date
			)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			is_verified = VALUES(is_verified)
	`

	binaryVerifyBookRequestID, err := uuid.UUID(verifyBookRequest.VerifyBookRequestID()).MarshalBinary()
	if err != nil {
		return err
	}
	binaryTranslatorID, err := uuid.UUID(verifyBookRequest.TranslatorID()).MarshalBinary()
	if err != nil {
		return err
	}
	binaryBookID, err := uuid.UUID(verifyBookRequest.BookID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query,
		binaryVerifyBookRequestID,
		binaryTranslatorID,
		binaryBookID,
		verifyBookRequest.IsVerified(),
		verifyBookRequest.SendDate(),
	)

	return err
}

func (repo *VerifyBookRequestRepository) Delete(verifyBookRequestID model.VerifyBookRequestID) error {
	const query = `DELETE FROM verify_book_request WHERE verify_book_request_id = ?`

	binaryVerifyBookRequestID, err := uuid.UUID(verifyBookRequestID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryVerifyBookRequestID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		return model.ErrVerifyBookRequestNotFound
	}

	return err
}

func (repo *VerifyBookRequestRepository) FindByID(verifyBookRequestID model.VerifyBookRequestID) (model.VerifyBookRequest, error) {
	const query = `
		SELECT
			translator_id,
			book_id,
			is_verified,
			send_date
		FROM verify_book_request
		WHERE verify_book_request_id = ?
`

	var verifyBookRequest sqlxVerifyBookRequest
	binaryVerifyBookRequestID, err := uuid.UUID(verifyBookRequestID).MarshalBinary()
	if err != nil {
		return model.VerifyBookRequest{}, err
	}

	err = repo.connection.Get(&verifyBookRequest, query, binaryVerifyBookRequestID)
	if errors.Is(err, sql.ErrNoRows) {
		return model.VerifyBookRequest{}, model.ErrVerifyBookRequestNotFound
	}
	if err != nil {
		return model.VerifyBookRequest{}, err
	}

	return model.NewVerifyBookRequest(
		verifyBookRequestID,
		model.UserID(verifyBookRequest.TranslatorID),
		verifyBookRequest.BookID,
		maybe.Just(verifyBookRequest.IsVerified),
		verifyBookRequest.SendDate,
	), nil
}

type sqlxVerifyBookRequest struct {
	TranslatorID uuid.UUID `db:"translator_id"`
	BookID       uuid.UUID `db:"book_id"`
	IsVerified   bool      `db:"is_verified"`
	SendDate     time.Time `db:"send_date"`
}
