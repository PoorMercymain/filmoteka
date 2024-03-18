package errors

import "errors"

var (
	ErrDuplicateInJSON       = errors.New("duplicate key found in JSON")
	ErrWrongMIME             = errors.New("wrong MIME type used")
	ErrWrongRequestWithJSON  = errors.New("something is wrong in request")
	ErrNothingProvidedInJSON = errors.New("nothing provided in JSON")
)
