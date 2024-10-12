package model

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"time"
)

type VerifyBookRequestID = uuid.UUID

type VerifyBookRequest struct {
	verifyBookRequestID VerifyBookRequestID
	translatorID        UserID
	bookID              BookID
	isVerified          maybe.Maybe[bool]
	sendDate            time.Time
}

func NewVerifyBookRequest(
	verifyBookRequestID VerifyBookRequestID,
	translatorID UserID,
	bookID BookID,
	isVerified maybe.Maybe[bool],
	sendDate time.Time,
) VerifyBookRequest {
	return VerifyBookRequest{
		verifyBookRequestID: verifyBookRequestID,
		translatorID:        translatorID,
		bookID:              bookID,
		isVerified:          isVerified,
		sendDate:            sendDate,
	}
}

func (verifyBookRequest *VerifyBookRequest) VerifyBookRequestID() VerifyBookRequestID {
	return verifyBookRequest.verifyBookRequestID
}

func (verifyBookRequest *VerifyBookRequest) TranslatorID() UserID {
	return verifyBookRequest.translatorID
}

func (verifyBookRequest *VerifyBookRequest) BookID() BookID {
	return verifyBookRequest.bookID
}

func (verifyBookRequest *VerifyBookRequest) IsVerified() maybe.Maybe[bool] {
	return verifyBookRequest.isVerified
}

func (verifyBookRequest *VerifyBookRequest) SendDate() time.Time {
	return verifyBookRequest.sendDate
}

func (verifyBookRequest *VerifyBookRequest) SetIsVerified(isVerified maybe.Maybe[bool]) {
	verifyBookRequest.isVerified = isVerified
}
