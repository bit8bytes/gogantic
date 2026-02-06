package main

import (
	"fmt"

	"github.com/bit8bytes/gogantic/input/chat"
	"github.com/bit8bytes/gogantic/llm"
)

func main() {
	chatPrompt := chat.New([]llm.Message{
		{Role: "system", Content: "You are a helpful assistant that translates {{.inputLanguage}} to {{.outputLanguage}}."},
		{Role: "user", Content: "{{.text}}"},
		{Role: "something", Content: "{{.something}}"},
	})

	data := map[string]string{
		"inputLanguage":  "English",
		"outputLanguage": "French",
		"text":           "I love programming.",
		"something":      "Hello, world!",
	}

	formattedMessages, err := chatPrompt.Format(data)
	if err != nil {
		panic(err) // This is not something we really want to do in production.
	}

	for _, message := range formattedMessages {
		fmt.Printf("[%s] %s\n", message.Role, message.Content)
	}
}
