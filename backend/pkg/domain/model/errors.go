package model

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrNotDeleteAdmin = errors.New("not delete admin")

	ErrBookNotFound = errors.New("book not found")

	ErrBookChapterNotFound = errors.New("book chapter not found")

	ErrVerifyBookRequestNotFound = errors.New("verify book request not found")
)
