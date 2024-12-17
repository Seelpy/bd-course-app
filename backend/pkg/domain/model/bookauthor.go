package model

type BookAuthor struct {
	bookID BookID
	author AuthorID
}

func NewBookAuthor(
	bookID BookID,
	authorID AuthorID,
) BookAuthor {
	return BookAuthor{
		bookID: bookID,
		author: authorID,
	}
}

func (bookAuthor *BookAuthor) BookID() BookID {
	return bookAuthor.bookID
}

func (bookAuthor *BookAuthor) Author() AuthorID {
	return bookAuthor.author
}
