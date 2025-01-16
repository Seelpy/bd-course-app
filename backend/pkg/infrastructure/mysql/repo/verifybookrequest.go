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

type verifyBookRequestRepository struct {
	connection *sqlx.DB
}

func NewVerifyBookRequestRepository(connection *sqlx.DB) *verifyBookRequestRepository {
	return &verifyBookRequestRepository{connection: connection}
}

func (repo *verifyBookRequestRepository) NextID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (repo *verifyBookRequestRepository) Store(verifyBookRequest model.VerifyBookRequest) error {
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

	var isVerified sql.NullBool
	if verified, ok := verifyBookRequest.IsVerified().Get(); ok {
		isVerified = sql.NullBool{Bool: verified, Valid: true}
	} else {
		isVerified = sql.NullBool{Valid: false}
	}

	_, err = repo.connection.Exec(query,
		binaryVerifyBookRequestID,
		binaryTranslatorID,
		binaryBookID,
		isVerified,
		verifyBookRequest.SendDate(),
	)

	return err
}

func (repo *verifyBookRequestRepository) Delete(verifyBookRequestID model.VerifyBookRequestID) error {
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

func (repo *verifyBookRequestRepository) FindByID(verifyBookRequestID model.VerifyBookRequestID) (model.VerifyBookRequest, error) {
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

	isVerified := maybe.Nothing[bool]()
	if verifyBookRequest.IsVerified.Valid {
		isVerified = maybe.Just(verifyBookRequest.IsVerified.Bool)
	}

	return model.NewVerifyBookRequest(
		verifyBookRequestID,
		model.UserID(verifyBookRequest.TranslatorID),
		verifyBookRequest.BookID,
		isVerified,
		verifyBookRequest.SendDate,
	), nil
}

type sqlxVerifyBookRequest struct {
	TranslatorID uuid.UUID    `db:"translator_id"`
	BookID       uuid.UUID    `db:"book_id"`
	IsVerified   sql.NullBool `db:"is_verified"`
	SendDate     time.Time    `db:"send_date"`
}
