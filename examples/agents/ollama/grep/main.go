package main

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/llms/ollama"
	"github.com/bit8bytes/gogantic/runner"
)

func main() {
	llm := ollama.New(ollama.Model{
		Model:   "gemma3n:e2b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
		Format:  "json",
	})

	tools := []agents.Tool{
		ListDir{},
	}

	task := "List all files in folder agents/"
	agent := agents.New(llm, tools)
	if err := agent.Task(task); err != nil {
		panic(err)
	}

	ctx := context.TODO()
	r := runner.New(agent, 10, true)
	if err := r.Run(ctx); err != nil {
		panic(err)
	}

	finalAnswer, err := agent.Answer()
	if err != nil {
		panic(err)
	}
	fmt.Println(finalAnswer)
}
