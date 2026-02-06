package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/input/chat"
	"github.com/bit8bytes/gogantic/llm"
	"github.com/bit8bytes/gogantic/llm/ollama"
	"github.com/bit8bytes/gogantic/output"
	"github.com/bit8bytes/gogantic/output/separator"
	"github.com/bit8bytes/gogantic/pipe"
)

func main() {
	chatPrompt := chat.New([]llm.Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant that translates {{.inputLanguage}} to {{.outputLanguage}}."},
		{
			Role:    "user",
			Content: "{{.text}}",
		},
	})

	data := map[string]any{
		"inputLanguage":  "English",
		"outputLanguage": "Spanish",
		"text":           "I love programming.",
	}

	messages, err := chatPrompt.Format(data)
	if err != nil {
		panic(err)
	}

	model := ollama.Model{
		Model:   "llama3:8b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
	}

	llm := ollama.New(model)
	var parser output.Parser[[]string] = &separator.Space{}

	ctx := context.TODO()
	pipe := pipe.New(messages, llm, parser)
	result, _ := pipe.Invoke(ctx)
	fmt.Println(*result)
}
