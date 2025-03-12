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
	chatPrompt, _ := input.NewChatPromptTemplate([]models.MessageContent{
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
	formattedMessages, err := chatPrompt.FormatMessages(data)
	if err != nil {
		panic(err) // Don't panic in production.
	}

	// Setup model with the wanted options.
	llama3_8b_model := ollama.OllamaModel{
		Model:   "llama3:8b",
		Options: ollama.ModelOptions{NumCtx: 4096},
		Stream:  false,
	}

	// Create a new client with the model
	ollamaClient := ollama.New(llama3_8b_model)
	// Generate the content using the model based on the formatted messages.
	generatedContent, _ := ollamaClient.GenerateContent(context.Background(), formattedMessages)
	// Do something with the generated content...
	fmt.Println(generatedContent.Result)
}
