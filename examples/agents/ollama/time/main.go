package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/core/models/llms/ollama"
)

type GetTime struct{}

func main() {
	wizardlm2_7b := ollama.OllamaModel{
		Model:   "mistral:latest",
		Options: ollama.ModelOptions{NumCtx: 4096},
		Stream:  false,
		Stop:    []string{"\nObservation", "Observation"},
	}
	llm := ollama.New(wizardlm2_7b)

	tools := map[string]agents.Tool{
		"GetTime": GetTime{},
	}

	timeAgent := agents.New(llm, tools)
	timeAgent.Task("What time is it?")

	ctx := context.TODO()
	executor := agents.NewExecutor(timeAgent, agents.WithShowMessages())
	executor.Run(ctx)

	finalAnswer, _ := timeAgent.GetFinalAnswer()
	fmt.Println(finalAnswer)
}

func (t GetTime) Name() string { return "GetTime" }

func (t GetTime) Call(ctx context.Context, input string) (string, error) {
	currentTime := time.Now()
	fmtCurrentTime := currentTime.Format("2006-01-02 3:04:05 PM")
	return fmtCurrentTime, nil
}
