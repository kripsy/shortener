package models

import "fmt"

type UniqueError struct {
	Text string
	Err  error
}

func NewUniqueError(fieldName string, err error) error {
	return &UniqueError{
		Text: fmt.Sprintf("%v already exists", fieldName),
		Err:  err,
	}
}

func (ue *UniqueError) Error() string {
	return ue.Err.Error()
}

type IsDeletedError struct {
	Text string
	Err  error
}

func NewIsDeletedError(shortURL string, err error) error {
	return &IsDeletedError{
		Text: fmt.Sprintf("%v is deleted", shortURL),
		Err:  err,
	}
}

func (ue *IsDeletedError) Error() string {
	return ue.Err.Error()
}
