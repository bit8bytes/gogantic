package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/inputs/chats"
	"github.com/bit8bytes/gogantic/llms"
	"github.com/bit8bytes/gogantic/llms/ollama"
)

func main() {
	// We use a chat prompt from the core/input
	prompt := chats.New([]llms.Message{
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
	messages, _ := prompt.Execute(data)

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
