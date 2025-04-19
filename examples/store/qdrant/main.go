package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bit8bytes/gogantic/core/embedder"
	"github.com/bit8bytes/gogantic/core/models/llms/ollama"
	"github.com/bit8bytes/gogantic/store/qdrant"
)

func main() {
	// Embedding size of ollama3:8b = 4096
	llama3_8b_model := ollama.OllamaModel{
		Model:   "llama3:8b",
		Options: ollama.ModelOptions{NumCtx: 4096},
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

	qdrt.UseCollection("something") // Use this to switch collections
	err := qdrt.GetCollectionHealth(context.Background())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	similarPoints, _ := qdrt.SimilaritySearch(ctx, "I don't like to hike", 2)
	fmt.Println("Similar points:", similarPoints)
}
