package embedder

import (
	"context"

	"github.com/bit8bytes/gogantic/core/models"
)

type Embedder struct {
	LLM models.Model
}

func New(llm models.Model) *Embedder {
	return &Embedder{
		LLM: llm,
	}
}

func (e *Embedder) EmbedQuery(ctx context.Context, query string) (models.EmbeddingResponse, error) {
	return e.LLM.GenerateEmbedding(ctx, query)
}
