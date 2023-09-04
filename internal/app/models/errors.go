package models

import "fmt"

// UniqueError represent type for uniqueness conflict error.
type UniqueError struct {
	Err  error
	Text string
}

// NewUniqueError return new uniqueness conflict error.
func NewUniqueError(fieldName string, err error) error {
	return &UniqueError{
		Text: fmt.Sprintf("%v already exists", fieldName),
		Err:  err,
	}

}

// Error to meet error intefrace by UniqueError.
func (ue *UniqueError) Error() string {
	return ue.Err.Error()
}

// IsDeletedError represent type of error, means object was removed.
type IsDeletedError struct {
	Err  error
	Text string
}

// NewIsDeletedError return new error, means object was removed.
func NewIsDeletedError(shortURL string, err error) error {
	return &IsDeletedError{
		Text: fmt.Sprintf("%v is deleted", shortURL),
		Err:  err,
	}
}

// Error to meet error intefrace by IsDeletedError.
func (ue *IsDeletedError) Error() string {
	return ue.Err.Error()
}
