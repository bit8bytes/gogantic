package ollama

import (
	"github.com/bit8bytes/gogantic/llms"
)

type Model struct {
	Model     string   `json:"model"`
	Endpoint  string   `json:"endpoint"`
	Options   Options  `json:"options"`
	Stream    bool     `json:"stream"`
	Format    string   `json:"format,omitempty"`
	KeepAlive int64    `json:"keepalive,omitempty"`
	Stop      []string `json:"stop"`
}

type Options struct {
	NumCtx      int     `json:"num_ctx"`
	Temperature float64 `json:"temperature"`
}

type Prompt struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	Options   Options
	Format    string `json:"format,omitempty"`
	KeepAlive int64  `json:"keepalive,omitempty"`
}

type ChatResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Message            llms.Message
	Done               bool  `json:"done"`
	TotalDuration      int64 `json:"total_duration"`
	LoadDuration       int64 `json:"load_duration"`
	PromptEvalCount    int   `json:"prompt_eval_count"`
	PromptEvalDuration int64 `json:"prompt_eval_duration"`
	EvalCount          int   `json:"eval_count"`
	EvalDuration       int64 `json:"eval_duration"`
}

type Chat struct {
	Model     string         `json:"model"`
	Messages  []llms.Message `json:"messages"`
	Options   Options
	Stream    bool   `json:"stream"`
	Format    string `json:"format,omitempty"`
	KeepAlive int64  `json:"keepalive,omitempty"`
}
