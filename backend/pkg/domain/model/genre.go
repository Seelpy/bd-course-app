package model

import "github.com/gofrs/uuid"

type GenreID = uuid.UUID

type Genre struct {
	id   GenreID
	name string
}

func NewGenre(
	id GenreID,
	name string,
) Genre {
	return Genre{
		id:   id,
		name: name,
	}
}

func (genre *Genre) ID() GenreID {
	return genre.id
}

func (genre *Genre) Name() string {
	return genre.name
}

func (genre *Genre) SetName(name string) {
	genre.name = name
}
