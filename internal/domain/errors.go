package domain

import "errors"

var (
	ErrNotFound = errors.New("resource not found")
	ErrNoRows   = errors.New("no rows affected")
)
