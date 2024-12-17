package model

type BookGenre struct {
	bookID  BookID
	genreID GenreID
}

func NewBookGenre(
	bookID BookID,
	genreID GenreID,
) BookGenre {
	return BookGenre{
		bookID:  bookID,
		genreID: genreID,
	}
}

func (bookGenre *BookGenre) BookID() BookID {
	return bookGenre.bookID
}

func (bookGenre *BookGenre) GenreID() GenreID {
	return bookGenre.genreID
}
