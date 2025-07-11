package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bit8bytes/gogantic/agent"
	"github.com/bit8bytes/gogantic/llm/ollama"
	"github.com/bit8bytes/gogantic/runner"
	"github.com/bit8bytes/gogantic/tool"
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

	tools := map[string]tool.Tool{
		"GetTime": GetTime{},
	}

	timeAgent := agent.New(llm, tools)
	timeAgent.Task("What time is it?")

	ctx := context.TODO()
	runner := runner.New(timeAgent, runner.WithShowMessages())
	runner.Run(ctx)

	finalAnswer, _ := timeAgent.GetFinalAnswer()
	fmt.Println(finalAnswer)
}

func (t GetTime) Name() string { return "GetTime" }

func (t GetTime) Schema() string { return `()` }

func (t GetTime) Call(ctx context.Context, input tool.Input) (tool.Output, error) {
	currentTime := time.Now()
	fmtCurrentTime := currentTime.Format("2006-01-02 3:04:05 PM")
	return tool.Output{Content: fmtCurrentTime}, nil
}
