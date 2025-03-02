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

type BlogArticle struct {
	Topic string `json:"topic"`
}

func main() {
	chatPrompt, _ := input.NewChatPromptTemplate([]models.MessageContent{
		{Role: "system", Content: "You are a helpful assistant that generates an exciting and engaging blog article. The user will give you the topic. Keep it short for now."},
		{Role: "user", Content: "Topic: {{.Topic}}"},
	})

	data := BlogArticle{
		Topic: "Why coders aren't funny.",
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

	ollamaClient := ollama.NewOllamaClient(llama3_8b_model)
	var parser output.OutputParser[map[string]string] = &output.MarkdownOutputParser{}

	ctx := context.TODO()
	pipe := pipe.NewPipe(messages, ollamaClient, parser)
	result, _ := pipe.Invoke(ctx)

	for heading, content := range result {
		// Do something with the sections. E.g. let other agents validate the content.
		fmt.Printf("\n--- %s ---\n%s\n", heading, content)
	}

}
