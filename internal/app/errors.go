package app

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("resource conflict")
	ErrInvalidCommand  = errors.New("invalid command")
	ErrUnauthorized    = errors.New("actor is not authorized")
	ErrCommandRejected = errors.New("command rejected")
	ErrUnavailable     = errors.New("service unavailable")
)
