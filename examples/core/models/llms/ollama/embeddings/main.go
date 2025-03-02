package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/core/input"
	"github.com/bit8bytes/gogantic/core/models/llms/ollama"
)

func main() {

	companyNamePrompt, _ := input.NewPromptTemplate("What is a good name for a company that makes {{.product}}?")

	data := map[string]any{"product": "coloful socks"}
	companyNameFormattedPrompt, _ := companyNamePrompt.Format(data)

	llama3_8b_model := ollama.OllamaModel{
		Model:   "llama3:8b",
		Options: ollama.ModelOptions{NumCtx: 4096},
		Stream:  false,
	}

	ollamaClient := ollama.NewOllamaClient(llama3_8b_model)
	generatedContent, _ := ollamaClient.GenerateEmbedding(context.Background(), companyNameFormattedPrompt)
	fmt.Println(len(generatedContent.Embedding))
	fmt.Println(generatedContent.Embedding[:3]) // First 3 bytes...
}
