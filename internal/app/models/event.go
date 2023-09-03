package models

import (
	"github.com/google/uuid"
)

// Event represent type for record of shortener.
type Event struct {
	ShortURL      string `json:"short_url,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	UUID          int    `json:"uuid,omitempty"`
	UserID        int    `json:"user_id,omitempty"`
	IsDeleted     bool   `json:"is_deleted,omitempty"`
}

// BatchURL represent slice of shortener records.
type BatchURL []Event

// NewEvent return new Event pointer.
func NewEvent(shortURL, originalURL string, userID int) *Event {
	e := &Event{
		UUID:        int(uuid.New().ID()),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
	}
	return e
}

// NewEventWithoutPointer return new Event.
// It's optimized version of NewEvent that use less memory because not using pointer.
func NewEventWithoutPointer(shortURL, originalURL string, userID int) Event {
	return Event{
		UUID:        int(uuid.New().ID()),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
	}
}
