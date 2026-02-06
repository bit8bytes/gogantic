package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/input/chat"
	"github.com/bit8bytes/gogantic/llm"
	"github.com/bit8bytes/gogantic/llm/ollama"
	"github.com/bit8bytes/gogantic/output"
	"github.com/bit8bytes/gogantic/output/json"
	"github.com/bit8bytes/gogantic/pipe"
)

// Please make sure to use the `json:"xxx"` behind the type.
// It is necessary to the JSON converter
type Translation struct {
	Text           string `json:"text"`
	InputLanguage  string `json:"inputLanguage"`
	OutputLanguage string `json:"outputLanguage"`
}

func main() {
	chatPrompt := chat.New([]llm.Message{
		{Role: "system", Content: "You are a helpful assistant that translates {{.InputLanguage}} to {{.OutputLanguage}}."},
		{Role: "user", Content: "{{.Text}}"},
	})

	data := Translation{
		Text:           "I love programming.",
		InputLanguage:  "English",
		OutputLanguage: "Spanish",
	}

	messages, err := chatPrompt.Format(data)
	if err != nil {
		panic(err)
	}

	llama3_8b_model := ollama.Model{
		Model:   "llama3:8b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
	}

	ollamaClient := ollama.New(llama3_8b_model)
	var parser output.Parser[Translation] = &json.Parser[Translation]{}

	pipe := pipe.New(messages, ollamaClient, parser)
	// Invoke is going to add the the parser instructions to the prompt.
	// The model generates the content and the parser parses the output then.
	result, _ := pipe.Invoke(context.Background())
	fmt.Println("Translate from", result.InputLanguage, " to ", result.OutputLanguage)
	fmt.Println("Result: ", result.Text)
}
