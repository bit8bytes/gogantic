package llms

import "github.com/bit8bytes/gogantic/inputs/roles"

type Message struct {
	Role    roles.Role `json:"role"`
	Content string     `json:"content"`
}

type ContentResponse struct {
	Result string
}

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

type StreamHandler func(content string, done bool) error
