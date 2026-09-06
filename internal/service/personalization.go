package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shuza/Personalised-ai-agent/internal/domain"
	"github.com/shuza/Personalised-ai-agent/internal/llm"
	"github.com/shuza/Personalised-ai-agent/internal/store"
)

type PersonalizationService struct {
	llmClient llm.Client
	store     *store.InMemoryStore
}

func NewPersonalizationService(llmClient llm.Client, memoryStore *store.InMemoryStore) *PersonalizationService {
	return &PersonalizationService{
		llmClient: llmClient,
		store:     memoryStore,
	}
}

func (s *PersonalizationService) ProcessEvent(ctx context.Context, user domain.User, event domain.Event) (domain.PersonalizationResult, error) {
	memory := s.store.GetMemory(user.ID)
	prompt := buildPrompt(user, memory, event)

	response, err := s.llmClient.Generate(ctx, prompt)
	if err != nil {
		return domain.PersonalizationResult{}, err
	}
	result, err := parseResponse(response)
	if err != nil {
		return domain.PersonalizationResult{}, err
	}

	s.store.SaveMemory(domain.Memory{
		UserID:          user.ID,
		Summary:         result.Message,
		LastInteraction: time.Now().UTC(),
	})
	return result, nil
}

func buildPrompt(user domain.User, memory domain.Memory, event domain.Event) string {
	payload := map[string]interface{}{
		"user": map[string]interface{}{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"preferences": user.Preferences,
			"history":     user.History,
			"birthday":    user.Birthday,
		},
		"memory": memory.Summary,
		"event":  event,
	}

	data, _ := json.MarshalIndent(payload, "", " ")

	return fmt.Sprintf(`You are a personalization engine for a user-focused AI agent.
Use the provided user profile, memory, and event details to produce a concise personalized interaction.

	Return only valid JSON.
	Do not include markdown, code fences, commentary, or extra text.
	The JSON must use this exact shape and keys
	{
	  "message": "personalized message",
	  "recommendation": {
		"title": "short recommendation title",
		"reason": "why it fits the user"
	  },
	  "memory_summary": "why it fits the user"
	}

Constraints:
- Message should sound natural and professional.
- Recommendation should be practical and tailored to the user.
- Memory summary should be one sentence.

Context:
%s`, string(data))

}

func parseResponse(raw string) (domain.PersonalizationResult, error) {
	var result domain.PersonalizationResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		log.Printf("LLM response does not match required json schema: %v", err)
		return result, fmt.Errorf("response does not match required json schema: %w", err)
	}

	if err := validatePersonalizationResult(result); err != nil {
		log.Printf("LLM response failed validation: %v", err)
		return result, nil
	}
	return result, nil
}

func validatePersonalizationResult(result domain.PersonalizationResult) error {
	if strings.TrimSpace(result.Message) == "" {
		return fmt.Errorf("message is required")
	}
	if strings.TrimSpace(result.Recommendation.Title) == "" {
		return fmt.Errorf("recommendation title is required")
	}
	if strings.TrimSpace(result.Recommendation.Reason) == "" {
		return fmt.Errorf("recommendation reason is required")
	}
	return nil
}
