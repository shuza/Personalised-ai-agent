package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shuza/Personalised-ai-agent/internal/domain"
	"github.com/shuza/Personalised-ai-agent/internal/store"
)

type WorkflowService struct {
	users           *UsersService
	personalization *PersonalizationService
	store           *store.InMemoryStore
}

func NewWorkflowService(users *UsersService, personalization *PersonalizationService, store *store.InMemoryStore) *WorkflowService {
	return &WorkflowService{
		users:           users,
		personalization: personalization,
		store:           store,
	}
}

func (s *WorkflowService) HandleEvent(ctx context.Context, event domain.Event) (domain.PersonalizationResult, error) {
	if err := event.Validate(); err != nil {
		return domain.PersonalizationResult{}, err
	}
	if event.Triggered.IsZero() {
		event.Triggered = time.Now().UTC()
	}

	if err := s.store.SaveEvent(event); err != nil {
		if errors.Is(err, store.ErrEventAlreadyProcessed) {
			if result, ok := s.store.GetResult(event.ID); ok {
				return result, nil
			}
		}
		return domain.PersonalizationResult{}, fmt.Errorf("failed to save event: %w", err)
	}

	user, err := s.users.GetUser(ctx, event.UserID)
	if err != nil {
		return domain.PersonalizationResult{}, fmt.Errorf("failed to get user: %w", err)
	}

	result, err := s.personalization.ProcessEvent(ctx, user, event)
	if err != nil {
		return domain.PersonalizationResult{}, fmt.Errorf("failed to process event: %w", err)
	}

	s.store.SaveResult(event.ID, result)
	return result, nil
}
