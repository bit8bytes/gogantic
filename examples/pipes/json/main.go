package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/inputs/chats"
	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
	"github.com/bit8bytes/gogantic/llms/ollama"
	"github.com/bit8bytes/gogantic/outputs/json"
	"github.com/bit8bytes/gogantic/pipes"
)

func main() {
	chat := chats.New([]llms.Message{
		{
			Role: roles.System,
			Content: `
You are a helpful assistant that translates {{.InputLanguage}} to {{.OutputLanguage}}.
Return only the result.
`,
		},
		{
			Role:    roles.User,
			Content: "{{.Text}}",
		},
	})

	// `json:"xxx"` is required for [jsonout] parser
	type translation struct {
		Text           string `json:"text"`
		InputLanguage  string `json:"inputLanguage"`
		OutputLanguage string `json:"outputLanguage"`
	}

	data := translation{
		Text:           "I love programming.",
		InputLanguage:  "en",
		OutputLanguage: "es",
	}

	messages, err := chat.Execute(data)
	if err != nil {
		panic(err)
	}

	model := ollama.Model{
		Model:   "gemma3n:e2b",
		Format:  json.Format,
		Options: ollama.Options{NumCtx: 4096},
	}

	client := ollama.New(model)

	parser := json.NewParser[translation]()

	pipe := pipes.New(messages, client, parser)
	// Invoke is going to add the the parser instructions to the prompt.
	// The model generates the content and the parser parses the output then.
	// The result will be of type [translation].
	result, _ := pipe.Invoke(context.Background())
	fmt.Println("Translate from", result.InputLanguage, "to", result.OutputLanguage)
	fmt.Println("Result:", result.Text)
}
