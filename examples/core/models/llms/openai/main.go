package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bit8bytes/gogantic/core/input"
	"github.com/bit8bytes/gogantic/core/models"
	"github.com/bit8bytes/gogantic/core/models/llms/openai"
)

func main() {

	chatPrompt, _ := input.NewChatPromptTemplate([]models.MessageContent{
		{Role: "system", Content: "You are a helpful assistant that translates {{.inputLanguage}} to {{.outputLanguage}}."},
		{Role: "user", Content: "{{.text}}"},
	})

	data := map[string]interface{}{
		"inputLanguage":  "English",
		"outputLanguage": "French",
		"text":           "I love programming.",
	}

	formattedMessages, err := chatPrompt.FormatMessages(data)
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("No OPENAI_API_KEY")
	}

	stream := false
	stop := []string{"\nObservation", "Observation"}

	gpt_model := openai.Model{
		Model:  "gpt-3.5-turbo",
		APIKey: apiKey,
		Stream: &stream,
		Stop:   &stop,
	}

	openAiClient := openai.New(gpt_model)
	generatedContent, err := openAiClient.GenerateContent(context.TODO(), formattedMessages)
	if err != nil {
		fmt.Println("Error generating content", err)
	}
	fmt.Println(generatedContent.Result)
}
