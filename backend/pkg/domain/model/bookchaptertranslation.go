package model

type BookChapterTranslation struct {
	bookChapterID BookChapterID
	translatorID  UserID
	text          string
}

func NewBookChapterTranslation(
	bookChapterID BookChapterID,
	translatorID UserID,
	text string,
) BookChapterTranslation {
	return BookChapterTranslation{
		bookChapterID: bookChapterID,
		translatorID:  translatorID,
		text:          text,
	}
}

func (bookChapterTranslation *BookChapterTranslation) BookChapterID() BookChapterID {
	return bookChapterTranslation.bookChapterID
}

func (bookChapterTranslation *BookChapterTranslation) TranslatorID() UserID {
	return bookChapterTranslation.translatorID
}

func (bookChapterTranslation *BookChapterTranslation) Text() string {
	return bookChapterTranslation.text
}

func (bookChapterTranslation *BookChapterTranslation) SetText(text string) {
	bookChapterTranslation.text = text
}
