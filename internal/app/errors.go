package app

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource conflict")
	ErrInvalidCommand = errors.New("invalid command")
	ErrUnauthorized   = errors.New("actor is not authorized")
	ErrUnavailable    = errors.New("service unavailable")
)
