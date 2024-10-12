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
