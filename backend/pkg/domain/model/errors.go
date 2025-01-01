package model

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrNotDeleteAdmin = errors.New("not delete admin")
)
