package domain

import "time"

type User struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Email       string            `json:"email"`
	Preferences map[string]string `json:"preferences"`
	History     []string          `json:"history"`
	Birthday    string            `json:"birthday"`
}

type Memory struct {
	UserID          string    `json:"user_id"`
	Summary         string    `json:"summary"`
	LastInteraction time.Time `json:"last_interaction"`
}

type Recommendation struct {
	Title  string `json:"title"`
	Reason string `json:"reason"`
}

type PersonalizationResult struct {
	Message        string         `json:"message"`
	Recommendation Recommendation `json:"recommendation"`
	MemorySummary  string         `json:"memory_summary"`
}
