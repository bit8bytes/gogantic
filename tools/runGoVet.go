package tools

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

// RunGoVet runs "go vet ./..." in the given directory and reports findings.
type RunGoVet struct{}

func (t RunGoVet) Name() string { return "RunGoVet" }

func (t RunGoVet) Description() string {
	return `Run "go vet ./..." in a Go module or package directory and return findings.
Input: absolute path to a directory containing Go source files.
Output: vet findings, or "ok" if no issues are found.`
}

func (t RunGoVet) Execute(ctx context.Context, input Input) (Output, error) {
	dir := strings.TrimSpace(input.Content)
	if dir == "" {
		return Output{Content: "error: provide an absolute path to a Go directory"}, nil
	}

	cmd := exec.CommandContext(ctx, "go", "vet", "./...")
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	cmd.Run() //nolint:errcheck // exit code is conveyed via output

	result := strings.TrimSpace(out.String())
	if result == "" {
		return Output{Content: "ok"}, nil
	}
	return Output{Content: result}, nil
}
