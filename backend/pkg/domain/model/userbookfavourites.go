package model

type UserBookFavouritesType int16

const (
	READING UserBookFavouritesType = iota
	PLANNED
	DEFERRED
	READ
	DROPPED
	FAVORITE
)

type UserBookFavourites struct {
	userID                 UserID
	bookID                 BookID
	userBookFavouritesType UserBookFavouritesType
}

func NewUserBookFavourites(
	userID UserID,
	bookID BookID,
	userBookFavouritesType UserBookFavouritesType,
) UserBookFavourites {
	return UserBookFavourites{
		userID:                 userID,
		bookID:                 bookID,
		userBookFavouritesType: userBookFavouritesType,
	}
}

func (userBookFavourites *UserBookFavourites) UserID() UserID {
	return userBookFavourites.userID
}

func (userBookFavourites *UserBookFavourites) BookID() BookID {
	return userBookFavourites.bookID
}

func (userBookFavourites *UserBookFavourites) UserBookFavouritesType() UserBookFavouritesType {
	return userBookFavourites.userBookFavouritesType
}

func (userBookFavourites *UserBookFavourites) SetUserBookFavouritesType(userBookFavouritesType UserBookFavouritesType) {
	userBookFavourites.userBookFavouritesType = userBookFavouritesType
}
