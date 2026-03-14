package tools

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

// RunGoBuild runs "go build ./..." in the given directory and reports compiler errors.
type RunGoBuild struct{}

func (t RunGoBuild) Name() string { return "RunGoBuild" }

func (t RunGoBuild) Description() string {
	return `Run "go build ./..." in a Go module or package directory and return compiler errors.
Input: absolute path to a directory containing Go source files.
Output: compiler errors, or "ok" if the build succeeds.`
}

func (t RunGoBuild) Execute(ctx context.Context, input Input) (Output, error) {
	dir := strings.TrimSpace(input.Content)
	if dir == "" {
		return Output{Content: "error: provide an absolute path to a Go directory"}, nil
	}

	cmd := exec.CommandContext(ctx, "go", "build", "./...")
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
