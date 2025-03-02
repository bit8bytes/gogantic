package openai

import (
	"github.com/bit8bytes/gogantic/core/models"
)

type Model struct {
	Model       string  `json:"model"`
	APIKey      string  `json:"api_key"`
	Temperature float32 `json:"temperature"`
}

type OpenAIRequest struct {
	Model       string                  `json:"model"`
	Messages    []models.MessageContent `json:"messages"`
	Temperature float32                 `json:"temperature"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   OpenAIResponseUsage
	Choices []OpenAIResponseChoice
}

type OpenAIResponseChoice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
}

type OpenAIResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
