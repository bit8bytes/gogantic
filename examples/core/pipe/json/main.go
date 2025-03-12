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

type Translation struct {
	Text           string `json:"text"`
	InputLanguage  string `json:"inputLanguage"`
	OutputLanguage string `json:"outputLanguage"`
}

func main() {
	chatPrompt, _ := input.NewChatPromptTemplate([]models.MessageContent{
		{Role: "system", Content: "You are a helpful assistant that translates {{.InputLanguage}} to {{.OutputLanguage}}."},
		{Role: "user", Content: "{{.Text}}"},
	})

	data := Translation{
		Text:           "I love programming.",
		InputLanguage:  "English",
		OutputLanguage: "Spanish",
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
	var parser output.OutputParser[Translation] = &output.JsonOutputParser[Translation]{}

	pipe := pipe.New(messages, ollamaClient, parser)
	// Invoke is going to add the the parser instructions to the prompt.
	// The model generates the content and the parser parses the output then.
	result, _ := pipe.Invoke(context.Background())
	fmt.Println("Translate from", result.InputLanguage, " to ", result.OutputLanguage)
	fmt.Println("Result: ", result.Text)
}
