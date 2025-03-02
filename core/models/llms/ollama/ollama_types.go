package ollama

import (
	"github.com/bit8bytes/gogantic/core/models"
)

type OllamaModel struct {
	Model     string       `json:"model"`
	Endpoint  string       `json:"endpoint"`
	Options   ModelOptions `json:"options"`
	Stream    bool         `json:"stream"`
	Format    string       `json:"format,omitempty"`
	KeepAlive int64        `json:"keepalive,omitempty"`
	Stop      []string     `json:"stop"`
}

type ModelOptions struct {
	NumCtx      int     `json:"num_ctx"`
	Temperature float64 `json:"temperature"`
}

type OllamaPromptRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	Options   ModelOptions
	Format    string `json:"format,omitempty"`
	KeepAlive int64  `json:"keepalive,omitempty"`
}

type ChatResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Message            models.MessageContent
	Done               bool  `json:"done"`
	TotalDuration      int64 `json:"total_duration"`
	LoadDuration       int64 `json:"load_duration"`
	PromptEvalCount    int   `json:"prompt_eval_count"`
	PromptEvalDuration int64 `json:"prompt_eval_duration"`
	EvalCount          int   `json:"eval_count"`
	EvalDuration       int64 `json:"eval_duration"`
}

type OllamaChatRequest struct {
	Model     string                  `json:"model"`
	Messages  []models.MessageContent `json:"messages"`
	Options   ModelOptions
	Stream    bool   `json:"stream"`
	Format    string `json:"format,omitempty"`
	KeepAlive int64  `json:"keepalive,omitempty"`
}
