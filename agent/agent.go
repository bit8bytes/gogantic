package agent

import (
	"context"
	"log/slog"
)

type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, input string) (string, error)
}

type Agent struct {
	tools  []Tool
	logger *slog.Logger
}

type Options struct {
	logger *slog.Logger
}

func WithLogger(logger *slog.Logger) func(*Options) {
	return func(o *Options) {
		o.logger = logger
	}
}

func New(tools []Tool, optsFunc ...func(*Options)) *Agent {
	opts := &Options{
		logger: slog.Default(),
	}

	for _, f := range optsFunc {
		f(opts)
	}

	return &Agent{
		tools:  tools,
		logger: opts.logger,
	}
}

func (a *Agent) Execute(ctx context.Context, toolName, input string) (string, error) {
	for _, tool := range a.tools {
		if tool.Name() == toolName {
			a.logger.Info("executing tool", "tool", toolName, "input", input)
			return tool.Execute(ctx, input)
		}
	}
	return "", ErrToolNotFound
}

func (a *Agent) ListTools() []string {
	var names []string
	for _, tool := range a.tools {
		names = append(names, tool.Name())
	}
	return names
}

func (a *Agent) GetToolDescription(toolName string) string {
	for _, tool := range a.tools {
		if tool.Name() == toolName {
			return tool.Description()
		}
	}
	return ""
}
