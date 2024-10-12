package model

type BookRating struct {
	bookID BookID
	userID UserID
	value  int
}
