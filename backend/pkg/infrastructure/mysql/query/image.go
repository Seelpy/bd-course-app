package query

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type ImageQueryService interface {
	FindByID(imageID model.ImageID) (string, error)
}

type imageQueryService struct {
	connection *sqlx.DB
}

func NewImageQueryService(connection *sqlx.DB) *imageQueryService {
	return &imageQueryService{connection: connection}
}

func (service *imageQueryService) FindByID(imageID model.ImageID) (string, error) {
	const query = `
		SELECT i.path
		FROM image i
		WHERE i.image_id = ?;
	`

	binaryImageID, err := uuid.UUID(imageID).MarshalBinary()
	if err != nil {
		return "", err
	}

	var imageData string
	err = service.connection.Get(&imageData, query, binaryImageID)
	if err != nil {
		return "", err
	}

	return imageData, nil
}
