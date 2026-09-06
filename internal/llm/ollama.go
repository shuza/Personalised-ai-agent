package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type OllamaClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type generateResponse struct {
	Response string `json:"response"`
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (c *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(generateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal ollama request: %w", err)
	}

	log.Printf("ollama request body: %s", string(body))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to send ollama request: %v", err)
		return "", fmt.Errorf("failed to send ollama request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Ollama request failed with status code %s", response.Status)
		payload, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return "", fmt.Errorf("ollama request failed with status code %s: %s", response.Status, strings.TrimSpace(string(payload)))
	}

	rawResponseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read ollama response body: %v", err)
		return "", fmt.Errorf("failed to read ollama response body: %w", err)
	}
	log.Printf("ollama response body: %s", string(rawResponseBody))

	var parsed generateResponse
	if err = json.Unmarshal(rawResponseBody, &parsed); err != nil {
		return "", fmt.Errorf("failed to decode ollama response: %w", err)
	}

	return strings.TrimSpace(parsed.Response), nil
}
