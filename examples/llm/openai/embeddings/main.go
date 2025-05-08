package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bit8bytes/gogantic/core/input"
	"github.com/bit8bytes/gogantic/core/models/llms/openai"
)

func main() {
	companyNamePrompt, _ := input.NewPromptTemplate("What is a good name for a company that makes {{.product}}?")

	data := map[string]any{"product": "coloful socks"}
	companyNameFormattedPrompt, _ := companyNamePrompt.Format(data)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("No OPENAI_API_KEY")
	}

	stream := false

	text_embedding_ada_002 := openai.Model{
		Model:  "text-embedding-ada-002",
		APIKey: apiKey,
		Stream: &stream,
	}

	openaiClient := openai.New(text_embedding_ada_002)
	generatedContent, err := openaiClient.GenerateEmbedding(
		context.Background(),
		companyNameFormattedPrompt,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(generatedContent.Embedding))
	fmt.Println(generatedContent.Embedding[:3]) // First 3 bytes...
}
