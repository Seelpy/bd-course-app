package service

import (
	"github.com/gofrs/uuid"
	"server/pkg/domain/model"
)

type ImageService interface {
	StoreImage(imageData string) (model.ImageID, error)
	DeleteImage(imageID model.ImageID) error
}

type imageService struct {
	imageRepo ImageRepository
}

func NewImageService(imageRepo ImageRepository) *imageService {
	return &imageService{
		imageRepo: imageRepo,
	}
}

type ImageRepository interface {
	NextID() uuid.UUID
	Store(image model.Image) error
	Delete(imageID model.ImageID) error
}

func (service *imageService) StoreImage(imageData string) (model.ImageID, error) {
	imageID := model.ImageID(service.imageRepo.NextID())

	image := model.NewImage(
		imageID,
		imageData,
	)

	return imageID, service.imageRepo.Store(image)
}

func (service *imageService) DeleteImage(imageID model.ImageID) error {
	return service.imageRepo.Delete(imageID)
}
