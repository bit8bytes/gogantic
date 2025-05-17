package qdrant

import (
	"context"
	"net/url"

	"github.com/bit8bytes/gogantic/embedder"
	"github.com/bit8bytes/gogantic/store"
)

type QdrantStore struct {
	Embedder   *embedder.Embedder
	Collection string
	Url        *url.URL
	ApiKey     string
}

type Options func(*QdrantStore)

func WithApiKey(apiKey string) Options {
	return func(qs *QdrantStore) {
		qs.ApiKey = apiKey
	}
}

func WithUrl(url *url.URL) Options {
	return func(qs *QdrantStore) {
		qs.Url = url
	}
}

func WithCollection(collection string) Options {
	return func(qs *QdrantStore) {
		qs.Collection = collection
	}
}

func WithEmbedder(embedder *embedder.Embedder) Options {
	return func(qs *QdrantStore) {
		qs.Embedder = embedder
	}
}

func New(opts ...Options) *QdrantStore {
	qs := &QdrantStore{}

	for _, opt := range opts {
		opt(qs)
	}

	if qs.Embedder == nil {
		panic("Embedder is required")
	}

	if qs.Url == nil {
		panic("URL is required")
	}

	return qs
}

func (qs *QdrantStore) UseCollection(collection string) {
	qs.Collection = collection
}

func (qs *QdrantStore) AddDocuments(ctx context.Context, docs []store.Document) (string, error) {
	upsertPointIds := qs.createUpsertPointIds(docs)
	metadatas := qs.createMetadatas(docs)
	contents := qs.createDocumentContent(docs)
	embeddedContents := qs.embedDocumentContents(ctx, contents)

	upsertPoints := UpsertPointsRequest{}
	upsertPoints.Batch.IDs = upsertPointIds
	upsertPoints.Batch.Payloads = metadatas
	upsertPoints.Batch.Vectors = embeddedContents

	response, err := qs.upsertPoints(ctx, upsertPoints)
	if err != nil {
		return response.Status, err
	}

	return response.Status, nil
}

func (qs *QdrantStore) SimilaritySearch(ctx context.Context, query string, limit int) ([]string, error) {
	// 1. Embed query
	embeddedResponse, err := qs.Embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	// 2. Search similar points in vectorstore
	payload := SearchPointsRequest{
		Limit:          limit,
		WithPayload:    true,
		WithVector:     false,
		Vector:         embeddedResponse.Embedding,
		ScoreThreshold: 0.2, // This may different from use case to use case.
	}

	similarPoints, err := qs.searchPoints(ctx, payload)
	if err != nil {
		return nil, err
	}

	// 3. Retrieve content field from payload
	var similarPointsContent []string
	for _, similarPoint := range similarPoints.Result {
		similarPointsContent = append(similarPointsContent, similarPoint.Payload["content"].(string))
	}

	return similarPointsContent, nil
}
