package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/input/prompt"
	"github.com/bit8bytes/gogantic/llm/ollama"
)

func main() {

	companyNamePrompt := prompt.New("What is a good name for a company that makes {{.product}}?")

	data := map[string]any{"product": "coloful socks"}
	companyNameFormattedPrompt, _ := companyNamePrompt.Format(data)

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
