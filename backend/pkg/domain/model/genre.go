package model

import "github.com/gofrs/uuid"

type GenreID = uuid.UUID

type Genre struct {
	id   GenreID
	name string
}
