package embedder

import (
	"context"

	"github.com/bit8bytes/gogantic/llm"
)

type Embedder struct {
	LLM llm.Model
}

func New(llm llm.Model) *Embedder {
	return &Embedder{
		LLM: llm,
	}
}

func (e *Embedder) EmbedQuery(ctx context.Context, query string) (llm.EmbeddingResponse, error) {
	return e.LLM.GenerateEmbedding(ctx, query)
}
