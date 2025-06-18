package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/url"

	"github.com/bit8bytes/gogantic/embedder"
	"github.com/bit8bytes/gogantic/llm/ollama"
	"github.com/bit8bytes/gogantic/store"
	"github.com/bit8bytes/gogantic/store/qdrant"
)

func main() {
	// Embedding size of ollama3:8b = 4096
	llama3_8b_model := ollama.Model{
		Model:   "llama3:8b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
	}

	ollamaClient := ollama.New(llama3_8b_model)
	ollamaEmbedder := embedder.New(ollamaClient)

	qdrantUrl := &url.URL{Scheme: "http", Host: "localhost:6333"}
	qdrt := qdrant.New(
		qdrant.WithCollection("something"),
		qdrant.WithEmbedder(ollamaEmbedder),
		qdrant.WithUrl(qdrantUrl),
	)

	qdrt.UseCollection("something")
	err := qdrt.GetCollectionHealth(context.Background())
	if err != nil {
		qdrt.CreateCollection(context.Background(), "something", qdrant.CreateCollectionRequest{
			Vectors: struct {
				Size     int    "json:\"size\""
				Distance string "json:\"distance\""
			}{
				Size:     4096,
				Distance: "Cosine",
			},
		})
	}

	docs := []store.Document{
		{
			ID:       fmt.Sprint(rand.Int()),
			Content:  "Take a leisurely walk in the park and enjoy the fresh air.",
			Metadata: map[string]any{"content": "Take a leisurely walk in the park and enjoy the fresh air."},
		},
		{
			ID:       fmt.Sprint(rand.Int()),
			Content:  "Visit a local museum and discover something new.",
			Metadata: map[string]any{"content": "Visit a local museum and discover something new."},
		},
	}

	ctx := context.Background()
	status, err := qdrt.AddDocuments(ctx, docs)
	if err != nil {
		fmt.Println("Something went wrong...", err)
	}
	fmt.Println("Everything seems to be: ", status)

	points, err := qdrt.CountPoints(context.Background())
	if err != nil {
		fmt.Println("Something went wrong...", err)
	}
	fmt.Println("Number of points in the collection:", points.Result.Count)
}
