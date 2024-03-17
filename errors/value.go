package errors

import "errors"

var (
	ErrUnknownGender         = errors.New("unknown gender used")
	ErrNoNameProvided        = errors.New("name not found in request or is empty")
	ErrNoIDProvided          = errors.New("id not found in request")
	ErrIDIsNotANumber        = errors.New("not a numeric id provided")
	ErrNoTitleProvided       = errors.New("title not found in request or is empty")
	ErrTitleTooLong          = errors.New("title is too long (150 characters is the limit)")
	ErrNoDescriptionProvided = errors.New("description not found in request or is empty")
	ErrDescriptionTooLong    = errors.New("description is too long (1000 characters is the limit)")
	ErrWrongRatingValue      = errors.New("rating should be in range [0, 10]")
)
