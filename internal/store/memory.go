package store

import (
	"errors"
	"sync"
	"time"

	"github.com/shuza/Personalised-ai-agent/internal/domain"
)

var ErrEventAlreadyProcessed = errors.New("event already processed")

type InMemoryStore struct {
	mu      sync.RWMutex
	users   map[string]domain.User
	memory  map[string]domain.Memory
	events  map[string]domain.Event
	results map[string]domain.PersonalizationResult
}

func NewInMemoryStore(defaultEmail string) *InMemoryStore {
	return &InMemoryStore{
		users: map[string]domain.User{
			"user-123": {
				ID:    "user-123",
				Name:  "Lily Mathew",
				Email: defaultEmail,
				Preferences: map[string]string{
					"favorite_cuisine": "japanese",
					"ambience":         "quite",
					"budget":           "50",
					"event_style":      "party",
					"channel":          "emails",
				},
				History: []string{
					"mentioned she prefers Japanese restaurants with calm atmosphere",
					"set dinner budget around 50",
					"love a rooftop birthday party event with live music",
				},
				Birthday: "1994-09-07",
			},
		},
		memory: map[string]domain.Memory{
			"user-123": {
				UserID:          "user-123",
				Summary:         "Prefers Japanese food in quite places, keeps a budget around 50, and enjoys party-style events.",
				LastInteraction: time.Now().Add(-48 * time.Hour),
			},
		},
		events:  make(map[string]domain.Event),
		results: make(map[string]domain.PersonalizationResult),
	}
}

func (s *InMemoryStore) GetUser(userID string) (domain.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[userID]
	return user, ok
}

func (s *InMemoryStore) GetMemory(userID string) domain.Memory {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.memory[userID]
}

func (s *InMemoryStore) SaveMemory(memory domain.Memory) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.memory[memory.UserID] = memory
}

func (s *InMemoryStore) SaveEvent(event domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; exists {
		return ErrEventAlreadyProcessed
	}
	return nil
}

func (s *InMemoryStore) SaveResult(eventID string, result domain.PersonalizationResult) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.results[eventID] = result
}

func (s *InMemoryStore) GetResult(eventID string) (domain.PersonalizationResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, ok := s.results[eventID]
	return result, ok
}
