package chat

import (
	"testing"

	"github.com/bit8bytes/gogantic/llm"
)

func TestNewChatPromptTemplate(t *testing.T) {
	messages := []llm.Message{
		{Role: "system", Content: "You are a helpful assistant that translates {{.inputLanguage}} to {{.outputLanguage}}."},
		{Role: "user", Content: "{{.text}}"},
	}

	chatPrompt, err := New(messages)
	if err != nil {
		t.Fatalf("unexpected error creating chat prompt template: %v", err)
	}

	if chatPrompt == nil {
		t.Errorf("expected chat prompt to be initialized, got nil")
	}
	if len(chatPrompt.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(chatPrompt.Messages))
	}
}

func TestChatPromptTemplateFormatMessages(t *testing.T) {
	messages := []llm.Message{
		{Role: "system", Content: "You are a helpful assistant that translates {{.inputLanguage}} to {{.outputLanguage}}."},
		{Role: "user", Content: "{{.text}}"},
	}

	chatPrompt, err := New(messages)
	if err != nil {
		t.Fatalf("unexpected error creating chat prompt template: %v", err)
	}

	data := map[string]interface{}{
		"inputLanguage":  "English",
		"outputLanguage": "French",
		"text":           "I love programming.",
	}

	formattedMessages, err := chatPrompt.Format(data)
	if err != nil {
		t.Fatalf("unexpected error formatting chat messages: %v", err)
	}

	expectedMessages := []llm.Message{
		{Role: "system", Content: "You are a helpful assistant that translates English to French."},
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
