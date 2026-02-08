package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/agents/tools"
	"github.com/bit8bytes/gogantic/llms/ollama"
	"github.com/bit8bytes/gogantic/runner"
)

type GetTime struct{}

func main() {
	model := ollama.Model{
		Model:   "gemma3n:e2b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
		Stop:    []string{"\nObservation", "Observation"},
	}
	llm := ollama.New(model)

	tools := map[string]agents.Tool{
		"GetTime": GetTime{},
	}

	timeAgent := agents.New(llm, tools)
	timeAgent.Task("What time is it?")

	ctx := context.TODO()
	runner := runner.New(timeAgent, runner.WithShowMessages())
	runner.Run(ctx)

	finalAnswer, _ := timeAgent.GetFinalAnswer()
	fmt.Println(finalAnswer)
}

func (t GetTime) Name() string { return "GetTime" }

func (t GetTime) Schema() string { return `()` }

func (t GetTime) Call(ctx context.Context, input tools.Input) (tools.Output, error) {
	currentTime := time.Now()
	fmtCurrentTime := currentTime.Format("2006-01-02 3:04:05 PM")
	return tools.Output{Content: fmtCurrentTime}, nil
}
