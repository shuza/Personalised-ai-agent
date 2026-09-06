package config

import "os"

type Config struct {
	HTTPAddress  string
	LLMProvider  string
	OllamaURL    string
	OllamaModel  string
	OpenAIModel  string
	OpenAIAPIKey string
	UserAPIBase  string
	DefaultEmail string
}

func Load() Config {
	return Config{
		HTTPAddress:  getEnv("HTTP_ADDRESS", ":8080"),
		LLMProvider:  getEnv("LLM_PROVIDER", "ollama"),
		OllamaURL:    getEnv("OLLAMA_URL", "http://localhost:11434"),
		OllamaModel:  getEnv("OLLAMA_MODEL", "llama3.1:8b"),
		OpenAIModel:  getEnv("OPENAI_MODEL", "gpt-4.1-mini"),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		UserAPIBase:  getEnv("USER_API_BASE", "https://localhost:8081"),
		DefaultEmail: getEnv("DEFAULT_EMAIL", "dummy@example.com"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
