package tools

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// ReadFile reads the full content of a Go source file.
type ReadFile struct{}

func (t ReadFile) Name() string { return "ReadFile" }

func (t ReadFile) Description() string {
	return `Read the full content of a file.
Input: absolute path to a file.
Output: file content as plain text.
Use this to read the complete source of a file when you need to review its full content.`
}

func (t ReadFile) Execute(ctx context.Context, input Input) (Output, error) {
	path := strings.TrimSpace(input.Content)
	if path == "" {
		return Output{Content: `{"error":"provide a file path"}`}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Output{Content: fmt.Sprintf(`{"error":%q}`, err.Error())}, nil
	}

	return Output{Content: string(data)}, nil
}
