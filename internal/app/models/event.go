package models

import (
	"github.com/google/uuid"
)

type Event struct {
	UUID          int    `json:"uuid,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	UserID        int    `json:"user_id,omitempty"`
}

type BatchURL []Event

func NewEvent(shortURL, originalURL string, userID int) *Event {
	e := &Event{
		UUID:        int(uuid.New().ID()),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
	}
	return e
}
