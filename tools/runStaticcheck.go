package tools

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

// RunStaticcheck runs "staticcheck ./..." in the given directory and reports findings.
type RunStaticcheck struct{}

func (t RunStaticcheck) Name() string { return "RunStaticcheck" }

func (t RunStaticcheck) Description() string {
	return `Run "staticcheck ./..." in a Go module or package directory and return findings.
Requires staticcheck to be installed (https://staticcheck.io).
Input: absolute path to a directory containing Go source files.
Output: staticcheck findings, "ok" if none, or an error if staticcheck is not installed.`
}

func (t RunStaticcheck) Execute(ctx context.Context, input Input) (Output, error) {
	dir := strings.TrimSpace(input.Content)
	if dir == "" {
		return Output{Content: "error: provide an absolute path to a Go directory"}, nil
	}

	if _, err := exec.LookPath("staticcheck"); err != nil {
		return Output{Content: "error: staticcheck not found in PATH — install it with: go install honnef.co/go/tools/cmd/staticcheck@latest"}, nil
	}

	cmd := exec.CommandContext(ctx, "staticcheck", "./...")
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
