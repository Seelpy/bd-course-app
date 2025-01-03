package service

import (
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"server/pkg/domain/model"
	"time"
)

type VerifyBookRequestService interface {
	CreateVerifyBookRequest(input CreateVerifyBookRequestInput) error
	AcceptVerifyBookRequest(input AcceptVerifyBookRequestInput) error
	DeleteVerifyBookRequest(input DeleteVerifyBookRequestInput) error
}

type verifyBookRequestService struct {
	verifyBookRequestRepo VerifyBookRequestRepository
}

func NewVerifyBookRequestService(
	verifyBookRequestRepo VerifyBookRequestRepository,
) *verifyBookRequestService {
	return &verifyBookRequestService{
		verifyBookRequestRepo: verifyBookRequestRepo,
	}
}

type VerifyBookRequestRepository interface {
	NextID() uuid.UUID
	Store(verifyBookRequest model.VerifyBookRequest) error
	Delete(verifyBookRequestID model.VerifyBookRequestID) error
	FindByID(verifyBookRequestID model.VerifyBookRequestID) (model.VerifyBookRequest, error)
}

type CreateVerifyBookRequestInput struct {
	TranslatorID model.UserID
	BookID       model.BookID
}

type DeleteVerifyBookRequestInput struct {
	VerifyBookRequestID model.VerifyBookRequestID
}

type AcceptVerifyBookRequestInput struct {
	VerifyBookRequestID model.VerifyBookRequestID
	Accept              bool
}

func (service *verifyBookRequestService) CreateVerifyBookRequest(input CreateVerifyBookRequestInput) error {
	verifyBookRequest := model.NewVerifyBookRequest(
		service.verifyBookRequestRepo.NextID(),
		input.TranslatorID,
		input.BookID,
		maybe.Nothing[bool](),
		time.Now(),
	)

	return service.verifyBookRequestRepo.Store(verifyBookRequest)
}

func (service *verifyBookRequestService) AcceptVerifyBookRequest(input AcceptVerifyBookRequestInput) error {
	verifyBookRequest, err := service.verifyBookRequestRepo.FindByID(input.VerifyBookRequestID)
	if err != nil {
		return err
	}

	verifyBookRequest.SetIsVerified(maybe.Just(input.Accept))

	return service.verifyBookRequestRepo.Store(verifyBookRequest)
}

func (service *verifyBookRequestService) DeleteVerifyBookRequest(input DeleteVerifyBookRequestInput) error {
	return service.verifyBookRequestRepo.Delete(input.VerifyBookRequestID)
}
