package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/core/input"
	"github.com/bit8bytes/gogantic/core/models"
	"github.com/bit8bytes/gogantic/core/models/llms/ollama"
	"github.com/bit8bytes/gogantic/core/output"
	"github.com/bit8bytes/gogantic/core/pipe"
)

func main() {
	chatPrompt, _ := input.NewChatPromptTemplate([]models.MessageContent{
		{Role: "system", Content: "You are a helpful assistant that translates {{.inputLanguage}} to {{.outputLanguage}}."},
		{Role: "user", Content: "{{.text}}"},
	})

	data := map[string]interface{}{
		"inputLanguage":  "English",
		"outputLanguage": "Spanish",
		"text":           "I love programming.",
	}

	messages, err := chatPrompt.FormatMessages(data)
	if err != nil {
		panic(err)
	}

	llama3_8b_model := ollama.OllamaModel{
		Model:   "llama3:8b",
		Options: ollama.ModelOptions{NumCtx: 4096},
		Stream:  false,
	}

	ollamaClient := ollama.New(llama3_8b_model)
	var parser output.OutputParser[[]string] = &output.SpaceSeparatedListOutputParser{}

	ctx := context.TODO()
	pipe := pipe.New(messages, ollamaClient, parser)
	result, _ := pipe.Invoke(ctx)
	fmt.Println(result)

}
