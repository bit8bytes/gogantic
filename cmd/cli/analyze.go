package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/llms/ollama"
	"github.com/bit8bytes/gogantic/runner"
	"github.com/bit8bytes/gogantic/stores/moderncsqlite"
	"github.com/bit8bytes/gogantic/tools"
)

func analyze(ctx context.Context, db *sql.DB, wd string) error {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	model := fs.String("model", "gemma3:12b", "Ollama model to use")
	numCtx := fs.Int("ctx", 16384, "Context window size")
	verbose := fs.Bool("v", false, "Print thought messages as the agent reasons")
	fs.Parse(os.Args[2:])

	targetPath := fs.Arg(0)
	if targetPath == "" {
		return fmt.Errorf("usage: beago analyze [flags] <path>")
	}

	storage, err := moderncsqlite.New(ctx, db)
	if err != nil {
		return fmt.Errorf("storage init: %w", err)
	}
	defer storage.Close()

	llm := ollama.New(ollama.Model{
		Model:   *model,
		Options: ollama.Options{NumCtx: *numCtx},
		Stream:  false,
		Format:  ollama.JSON,
	})

	goTools := []agents.Tool{
		tools.ReadFile{},
		tools.RunGoBuild{},
		tools.RunGoVet{},
	}

	agent, err := agents.NewReAct(ctx, llm, goTools, storage)
	if err != nil {
		return fmt.Errorf("agent init: %w", err)
	}

	task := buildTask(wd, targetPath)

	if err := agent.Task(ctx, task); err != nil {
		return fmt.Errorf("task: %w", err)
	}

	r := runner.New(agent, *verbose)
	if err := r.Run(ctx); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	finalAnswer, err := agent.Answer()
	if errors.Is(err, agents.ErrNoFinalAnswer) {
		fmt.Println("No final answer found")
		return nil
	}
	if err != nil {
		return fmt.Errorf("answer: %w", err)
	}

	fmt.Println(finalAnswer)
	return nil
}

func buildTask(wd, targetPath string) string {
	cfg, _ := load(wd)
	layer := cfg.MatchLayer(wd, targetPath)

	var sb strings.Builder

	fmt.Fprintf(&sb, `You are a Go code reviewer helping a junior developer fix issues before submitting a PR.

Project root: %s
Target file: %s

Steps:
1. Call ReadFile on the target file.
2. Call RunGoBuild on the project root to catch compilation errors.
3. Call RunGoVet on the project root to catch common bugs.
4. Call RunLinter on the project root to catch style and correctness issues.
`, wd, targetPath)

	if layer != nil {
		fmt.Fprintf(&sb, `5. Check the file content from step 1 against these architectural rules for the "%s" layer:

   - %s

6. Report all findings in the format below.
`, layer.Name, strings.Join(layer.Rules, "\n   - "))
	} else {
		sb.WriteString("5. Report all findings in the format below.\n")
	}

	sb.WriteString(`
Explain each issue clearly so a junior developer understands what to fix and why.

Format your final answer exactly like this:

REVIEW <target file>

No issues found.

Or if there are issues:

REVIEW <target file>

ISSUE: <short title>
  File:     <file>:<line>
  Rule:     <rule name or "style">
  Found:    <what the code does>
  Expected: <what it should do>
  Fix:      <concrete suggestion>

Only report real findings from the tools. No invented issues.`)

	return sb.String()
}
