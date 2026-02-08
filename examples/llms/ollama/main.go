package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/inputs/chats"
	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
	"github.com/bit8bytes/gogantic/llms/ollama"
)

func main() {
	messages := []llms.Message{
		{
			Role:    roles.System,
			Content: "Translate {{.InputLanguage}} to {{.OutputLanguage}}.",
		},
		{
			Role:    roles.System,
			Content: "Return only the concrete translation.",
		},
		{
			Role:    roles.User,
			Content: "{{.Text}}",
		},
	}

	type chatData struct {
		InputLanguage  string
		OutputLanguage string
		Text           string
	}

	data := chatData{
		InputLanguage:  "English",
		OutputLanguage: "French",
		Text:           "I love programming.",
	}

	formattedMessages, err := chats.New(messages).Execute(data)
	if err != nil {
		panic(err)
	}

	client := ollama.New(ollama.Model{
		Model:   "gemma3n:e2b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
	})
	generatedContent, _ := client.Generate(context.Background(), formattedMessages)
	fmt.Println(generatedContent.Result)
}
