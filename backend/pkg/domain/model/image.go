package model

import "github.com/gofrs/uuid"

type ImageID uuid.UUID

type Image struct {
	id   ImageID
	path string
}

func NewImage(
	id ImageID,
	path string,
) Image {
	return Image{
		id:   id,
		path: path,
	}
}

func (image *Image) ID() ImageID {
	return image.id
}

func (image *Image) Path() string {
	return image.path
}
