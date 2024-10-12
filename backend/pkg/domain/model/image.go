package model

import "github.com/gofrs/uuid"

type ImageID uuid.UUID

type Image struct {
	id   ImageID
	path string
}
