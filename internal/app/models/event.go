package models

import "github.com/google/uuid"

type Event struct {
	UUID          int    `json:"uuid,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
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
