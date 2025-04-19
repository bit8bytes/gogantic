package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"maps"

	"github.com/bit8bytes/gogantic/store"
	"github.com/google/uuid"
)

func (qs *QdrantStore) sendHttpRequestWithContext(ctx context.Context, method, url string, payload any, response any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+qs.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("unexpected status code: %d %s. Response: %s", resp.StatusCode, http.StatusText(resp.StatusCode), bodyString)
	}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func (qs *QdrantStore) createDocumentContent(docs []store.Document) []string {
	contents := make([]string, 0, len(docs))
	for _, doc := range docs {
		contents = append(contents, doc.Content)
	}
	return contents
}

func (qs *QdrantStore) embedDocumentContents(ctx context.Context, contents []string) [][]float32 {
	vectors := make([][]float32, 0, len(contents))
	for _, text := range contents {
		vector, _ := qs.Embedder.EmbedQuery(ctx, text)
		vectors = append(vectors, vector.Embedding)
	}
	return vectors
}

func (qs *QdrantStore) createMetadatas(docs []store.Document) []map[string]any {
	metadatas := make([]map[string]any, 0, len(docs))
	for i := range docs {
		metadata := make(map[string]any, len(docs[i].Metadata))
		maps.Copy(metadata, docs[i].Metadata)
		metadatas = append(metadatas, metadata)
	}
	return metadatas
}

func (qs *QdrantStore) createUpsertPointIds(docs []store.Document) []ID {
	ids := make([]ID, len(docs))
	for i := range ids {
		ids[i] = uuid.NewString()
	}
	return ids
}
