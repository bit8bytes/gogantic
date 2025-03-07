package agents

import (
	"context"
	"strings"
)

type Tool interface {
	Name() string
	Call(ctx context.Context, input string) (string, error)
}

func getToolNames(tools map[string]Tool) string {
	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		names = append(names, tool.Name())
	}
	return strings.Join(names, ", ")
}
