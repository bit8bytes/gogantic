package chats

import (
	"testing"

	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
)

func TestNewChatPromptTemplate(t *testing.T) {
	messages := []llms.Message{
		{
			Role:    "system",
			Content: "Translate {{.inputLanguage}} to {{.outputLanguage}}.",
		},
		{
			Role:    "user",
			Content: "{{.text}}",
		},
	}

	chatPrompt := New(messages)

	if chatPrompt == nil {
		t.Errorf("expected chat prompt to be initialized, got nil")
	}
	if len(chatPrompt.messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(chatPrompt.messages))
	}
}

func TestChatPromptTemplateFormatMessages(t *testing.T) {
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

	chatPrompt := New(messages)

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

	formattedMessages, err := chatPrompt.Execute(data)
	if err != nil {
		t.Fatalf("unexpected error formatting chat messages: %v", err)
	}

	expectedMessages := []llms.Message{
		{Role: "system", Content: "Translate English to French."},
		{Role: "user", Content: "I love programming."},
	}

	if len(formattedMessages) != len(expectedMessages) {
		t.Fatalf("expected %d formatted messages, got %d", len(expectedMessages), len(formattedMessages))
	}

	for i, msg := range formattedMessages {
		if msg.Role != expectedMessages[i].Role {
			t.Errorf("expected role %q, got %q", expectedMessages[i].Role, msg.Role)
		}
		if msg.Content != expectedMessages[i].Content {
			t.Errorf("expected content %q, got %q", expectedMessages[i].Content, msg.Content)
		}
	}
}
