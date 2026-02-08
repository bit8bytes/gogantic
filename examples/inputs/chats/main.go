package main

import (
	"fmt"

	"github.com/bit8bytes/gogantic/inputs/chats"
	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
)

func main() {
	messages := []llms.Message{
		{
			Role:    roles.System,
			Content: "Translate {{.InputLanguage}} to {{.OutputLanguage}}.",
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

	for _, message := range formattedMessages {
		fmt.Printf("[%s] %s\n", message.Role, message.Content)
	}
}
