package store

import "context"

type Map = map[string]any

type Document struct {
	Content  string
	Metadata Map
	Score    float32
}

type VectorStore interface {
	AddDocuments(ctx context.Context, docs []Document) ([]string, error)
	SimilaritySearch(ctx context.Context, query string, limit int) ([]string, error)
}
