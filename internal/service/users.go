package service

import (
	"context"
	"fmt"

	"github.com/shuza/Personalised-ai-agent/internal/domain"
	"github.com/shuza/Personalised-ai-agent/internal/store"
)

type UsersService struct {
	store *store.InMemoryStore
}

func NewUserService(memoryStore *store.InMemoryStore) *UsersService {
	return &UsersService{
		store: memoryStore,
	}
}

func (s *UsersService) GetUser(_ context.Context, userID string) (domain.User, error) {
	user, ok := s.store.GetUser(userID)
	if !ok {
		return domain.User{}, fmt.Errorf("user %s not found", userID)
	}
	return user, nil
}
