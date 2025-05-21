package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/input/chat"
	"github.com/bit8bytes/gogantic/llm"
	"github.com/bit8bytes/gogantic/llm/ollama"
	"github.com/bit8bytes/gogantic/output"
	"github.com/bit8bytes/gogantic/output/markdown"
	"github.com/bit8bytes/gogantic/pipe"
)

type BlogArticle struct {
	Topic string `json:"topic"`
}

func main() {
	chatPrompt, _ := chat.New([]llm.Message{
		{Role: "system", Content: "You are a helpful assistant that generates an exciting and engaging blog article. The user will give you the topic. Keep it short for now."},
		{Role: "user", Content: "Topic: {{.Topic}}"},
	})

	data := BlogArticle{
		Topic: "Why coders aren't funny.",
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
	var parser output.Parser[map[string]string] = &markdown.Parser{}

	ctx := context.TODO()
	pipe := pipe.New(messages, ollamaClient, parser)
	result, _ := pipe.Invoke(ctx)

	for heading, content := range *result {
		// Do something with the sections. E.g. let other agents validate the content.
		fmt.Printf("\n--- %s ---\n%s\n", heading, content)
	}

}
