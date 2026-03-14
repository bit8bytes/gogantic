package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"

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

	// These goTools are specifically designed for Golang.
	goTools := []agents.Tool{
		tools.GitDiff{},
		tools.ListDeclarations{},
		tools.FindCalls{},
		tools.FindUsages{},
		tools.GetFunctionSignature{},
		tools.GetStructFields{},
		tools.RunGoBuild{},
		tools.RunGoVet{},
	}

	agent, err := agents.NewReAct(ctx, llm, goTools, storage)
	if err != nil {
		return fmt.Errorf("agent init: %w", err)
	}

	task := fmt.Sprintf(`You are a Go code reviewer. Your job is to check if newly written code follows the same style and best practices as the rest of this codebase.

Project root: %s
Target file: %s

Follow these steps:
1. Call GitDiff with the target file path to get the exact changes.
2. Call ListDeclarations with the project root to find comparable existing code.
3. Call GetFunctionSignature on functions from the diff and on similar functions found in step 2.
4. Compare them. Report ONLY concrete deviations in: error handling style, naming conventions, doc comments, receiver naming, function signature shape.

Format your final answer exactly like this:

REVIEW <target file>

No issues found.

Or if there are deviations:

REVIEW <target file>

ISSUE: <short title>
  File:     <file>:<line>
  Expected: <what the codebase does>
  Found:    <what the new code does>

Report only concrete deviations. No prose outside this format.`, wd, targetPath)

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
