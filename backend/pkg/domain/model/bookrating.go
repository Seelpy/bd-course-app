package model

type BookRating struct {
	bookID BookID
	userID UserID
	value  int
}

func NewBookRating(
	bookID BookID,
	userID UserID,
	value int,
) BookRating {
	return BookRating{
		bookID: bookID,
		userID: userID,
		value:  value,
	}
}

func (bookRating *BookRating) BookID() BookID {
	return bookRating.bookID
}

func (bookRating *BookRating) UserID() UserID {
	return bookRating.userID
}

func (bookRating *BookRating) Value() int {
	return bookRating.value
}
