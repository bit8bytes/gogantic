package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/inputs/prompts"
	"github.com/bit8bytes/gogantic/llms/ollama"
)

func main() {

	companyNamePrompt := prompts.New("What is a good name for a company that makes {{.product}}?")

	data := map[string]any{"product": "coloful socks"}
	companyNameFormattedPrompt, _ := companyNamePrompt.Execute(data)

	llama3_8b_model := ollama.Model{
		Model:   "llama3:8b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
	}

	ollamaClient := ollama.New(llama3_8b_model)
	generatedContent, _ := ollamaClient.GenerateEmbedding(context.Background(), companyNameFormattedPrompt)
	fmt.Println(len(generatedContent.Embedding))
	fmt.Println(generatedContent.Embedding[:3]) // First 3 bytes...
}
