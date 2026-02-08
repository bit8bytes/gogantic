package embedder

import (
	"context"

	"github.com/bit8bytes/gogantic/llms"
)

type llm interface {
	GenerateEmbedding(ctx context.Context, prompt string) (llms.EmbeddingResponse, error)
}

type embedder struct {
	llm llm
}

func New(llm llm) *embedder {
	return &embedder{
		llm: llm,
	}
}

func (e *embedder) Embed(ctx context.Context, query string) (llms.EmbeddingResponse, error) {
	return e.llm.GenerateEmbedding(ctx, query)
}
