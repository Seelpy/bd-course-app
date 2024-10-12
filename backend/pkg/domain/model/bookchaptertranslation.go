package model

type BookChapterTranslation struct {
	bookChapterID BookChapterID
	translatorID  UserID
	text          string
}
