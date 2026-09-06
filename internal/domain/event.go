package domain

import (
	"errors"
	"time"
)

type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	UserID    string                 `json:"user_id"`
	Triggered time.Time              `json:"triggered"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func (e *Event) Validate() error {
	if e.ID == "" {
		return errors.New("event id is required")
	}
	if e.UserID == "" {
		return errors.New("user_id is required")
	}
	if e.Type == "" {
		return errors.New("event type is required")
	}
	return nil
}
