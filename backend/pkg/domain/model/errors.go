package model

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrNotDeleteAdmin = errors.New("not delete admin")

	ErrBookNotFound = errors.New("book not found")

	ErrBookChapterNotFound = errors.New("book chapter not found")

	ErrBookChapterTranslationNotFound = errors.New("book chapter translation not found")

	ErrVerifyBookRequestNotFound = errors.New("verify book request not found")

	ErrBookRatingNotFound = errors.New("book rating not found")

	ErrSessionReadingNotFound = errors.New("session reading not found")
)
