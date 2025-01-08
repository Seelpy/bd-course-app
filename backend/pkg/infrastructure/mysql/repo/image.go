package repo

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"server/pkg/domain/model"
)

type imageRepository struct {
	connection *sqlx.DB
}

func NewImageRepository(connection *sqlx.DB) *imageRepository {
	return &imageRepository{connection}
}

func (repo *imageRepository) NextID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (repo *imageRepository) Store(image model.Image) error {
	const query = `
		INSERT INTO
			image (
			      image_id,
			      path
			)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			path = VALUES(path)
	`

	binaryImageID, err := uuid.UUID(image.ID()).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.connection.Exec(query,
		binaryImageID,
		image.Path(),
	)

	return err
}

func (repo *imageRepository) Delete(imageID model.ImageID) error {
	const query = `DELETE FROM image WHERE image_id = ?`

	binaryImageID, err := uuid.UUID(imageID).MarshalBinary()
	if err != nil {
		return err
	}

	result, err := repo.connection.Exec(query, binaryImageID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		return model.ErrImageNotFound
	}

	return err
}
