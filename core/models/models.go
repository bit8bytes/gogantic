package models

import (
	"context"
)

type MessageContent struct {
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
	GenerateContent(ctx context.Context, messages []MessageContent) (*ContentResponse, error)
	GenerateEmbedding(ctx context.Context, prompt string) (EmbeddingResponse, error)
	StreamContent(ctx context.Context, messages []MessageContent, streamHandler StreamHandler) error
}
