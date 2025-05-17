package embedder

import (
	"context"

	"github.com/bit8bytes/gogantic/llm"
)

type Embedder struct {
	LLM llm.LLM
}

func New(llm llm.LLM) *Embedder {
	return &Embedder{
		LLM: llm,
	}
}

func (e *Embedder) EmbedQuery(ctx context.Context, query string) (llm.EmbeddingResponse, error) {
	return e.LLM.GenerateEmbedding(ctx, query)
}
