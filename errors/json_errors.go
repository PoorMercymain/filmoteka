package errors

import "errors"

var (
	ErrDuplicateInJSON = errors.New("duplicate key found in JSON")
	ErrWrongMIME       = errors.New("wrong MIME type used")
)
