package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/shuza/Personalised-ai-agent/internal/config"
	"github.com/shuza/Personalised-ai-agent/internal/llm"
	"github.com/shuza/Personalised-ai-agent/internal/service"
	"github.com/shuza/Personalised-ai-agent/internal/store"
	ttptransport "github.com/shuza/Personalised-ai-agent/internal/transport/http"
)

type App struct {
	router http.Handler
}

func New(cfg config.Config) *App {
	memoryStore := store.NewInMemoryStore(cfg.DefaultEmail)
	llmClient, err := newLLMClient(cfg)
	if err != nil {
		panic(err)
	}
	userService := service.NewUserService(memoryStore)
	personalizationService := service.NewPersonalizationService(llmClient, memoryStore)
	workflowService := service.NewWorkflowService(userService, personalizationService, memoryStore)
	handler := ttptransport.NewHandler(userService, workflowService)

	return &App{
		router: handler.Router(),
	}
}

func (a *App) Router() http.Handler {
	return a.router
}

func newLLMClient(cfg config.Config) (llm.Client, error) {
	switch strings.ToLower(cfg.LLMProvider) {
	case "ollama":
		return llm.NewOllamaClient(cfg.OllamaURL, cfg.OllamaModel), nil
	case "openai":
		return nil, fmt.Errorf("llm provider %q requires implementation before use", cfg.LLMProvider)
	default:
		return nil, fmt.Errorf("llm provider %q not supported", cfg.LLMProvider)
	}
}
