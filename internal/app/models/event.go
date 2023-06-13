package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Event struct {
	UUID          int    `json:"uuid,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
}

type BatchURL []Event

func NewEvent(shortURL string, originalURL string) *Event {
	e := &Event{
		UUID:        int(uuid.New().ID()),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
	return e
}

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
