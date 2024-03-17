package errors

import "errors"

var (
	ErrNotFoundInDB                  = errors.New("the requested entity does not exist in database")
	ErrActorNotBornBeforeFilmRelease = errors.New("one or more actors are not born before film release")
	ErrActorDoesNotExist             = errors.New("one or more actors mentioned in request does not exist in database")
)
