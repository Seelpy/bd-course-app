package model

import "github.com/gofrs/uuid"

type BookID = uuid.UUID

type Book struct {
	id          BookID
	coverID     ImageID
	title       string
	description string
	isPublished bool
}
