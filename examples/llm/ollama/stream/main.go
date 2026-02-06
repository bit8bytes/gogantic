package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/input/chat"
	"github.com/bit8bytes/gogantic/llm"
	"github.com/bit8bytes/gogantic/llm/ollama"
)

func main() {
	// We use a chat prompt from the core/input
	prompt := chat.New([]llm.Message{
		{Role: "system", Content: "You are a helpful assistant that translates {{.inputLanguage}} to {{.outputLanguage}}."},
		{Role: "user", Content: "{{.text}}"},
	})

	// Setup values for variables
	data := map[string]string{
		"inputLanguage":  "English",
		"outputLanguage": "French",
		"text":           "I love programming.",
	}

	// Prepare the chat prompt for the model
	messages, _ := prompt.Format(data)

	// Setup model with the wanted options.
	model := ollama.Model{
		Model:   "llama3:8b",
		Options: ollama.Options{NumCtx: 4096},
	}

	client := ollama.New(model)
	client.StreamContent(context.Background(), messages, func(content string, done bool) error {
		// Do something with the generated content...
		fmt.Print(content)
		return nil
	})
}
