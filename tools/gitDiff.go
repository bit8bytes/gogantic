package tools

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

// GitDiff runs "git diff HEAD -- <file>" for a given file path and returns the diff output.
type GitDiff struct{}

func (t GitDiff) Name() string { return "GitDiff" }

func (t GitDiff) Description() string {
	return `Run "git diff HEAD -- <file>" for a given file and return the unified diff.
Falls back to "git diff --cached -- <file>" for staged but uncommitted files.
Input: absolute path to a Go source file.
Output: unified diff showing what changed, or "no diff" if nothing changed.`
}

func (t GitDiff) Execute(ctx context.Context, input Input) (Output, error) {
	path := strings.TrimSpace(input.Content)
	if path == "" {
		return Output{Content: "error: provide an absolute path to a Go source file"}, nil
	}

	run := func(args ...string) string {
		var out bytes.Buffer
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Stdout = &out
		cmd.Stderr = &out
		cmd.Run() //nolint:errcheck // exit code is conveyed via output
		return strings.TrimSpace(out.String())
	}

	result := run("diff", "HEAD", "--", path)
	if result == "" {
		result = run("diff", "--cached", "--", path)
	}
	if result == "" {
		return Output{Content: "no diff"}, nil
	}
	return Output{Content: result}, nil
}
