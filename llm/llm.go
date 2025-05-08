package llm

import (
	"context"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ContentResponse struct {
	Result string
}

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

type StreamHandler func(content string, done bool) error

type Model interface {
	GenerateContent(ctx context.Context, messages []Message) (*ContentResponse, error)
	GenerateEmbedding(ctx context.Context, prompt string) (EmbeddingResponse, error)
	StreamContent(ctx context.Context, messages []Message, streamHandler StreamHandler) error
}
