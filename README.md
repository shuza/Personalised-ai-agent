# Personalized Agent for Event Celebration

This repository contains the code for a personalized agent designed to assist with event celebrations. The agent utilizes natural language processing and machine learning techniques to provide tailored recommendations and support for various event-related tasks.

## Features
- Event-driven personalization workflow
- User lookup
- Ollama-backed LLM generation
- Provider switch with reserved OpenAI configuration
- In-memory user context and memory store for development

## Run locally
1. Start Ollama and pull model:
```bash
ollama pull llama3.1:8b
ollama serve 
```
2. Update environment variables in `.env`
3. Run the API:
```bash
go run ./cmd/api/main.go 
```

# Endpoints
- `GET /healthz`
- `GET /v1/users/user-123`
- `POST /v1/events`
Sample event:
```json
{
  "id": "evt-birthday-001",
  "type": "birthday",
  "user_id": "user-123",
  "triggered_at": "2026-09-05T18:00:00Z",
  "metadata": {
    "source": "scheduler"
  }
}
```