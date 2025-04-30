package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/core/input"
	"github.com/bit8bytes/gogantic/core/models"
	"github.com/bit8bytes/gogantic/core/models/llms/ollama"
)

func main() {
	// We use a chat prompt from the core/input
	prompt, _ := input.NewChatPromptTemplate([]models.MessageContent{
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
	messages, _ := prompt.FormatMessages(data)

	// Setup model with the wanted options.
	model := ollama.OllamaModel{
		Model:   "llama3:8b",
		Options: ollama.ModelOptions{NumCtx: 4096},
	}

	client := ollama.New(model)
	client.StreamContent(context.Background(), messages, func(content string, done bool) error {
		// Do something with the generated content...
		fmt.Print(content)
		return nil
	})
}
