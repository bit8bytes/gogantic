package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/core/models/llms/ollama"
)

type TimeAgent struct{}

func main() {
	mistral_latest := ollama.OllamaModel{
		Model:   "mistral:latest",
		Options: ollama.ModelOptions{NumCtx: 4096},
		Stream:  false,
		Stop:    []string{"\nObservation", "Observation"},
	}
	llm := ollama.NewOllamaClient(mistral_latest)

	tools := map[string]agents.Tool{
		"TimeAgent": TimeAgent{},
	}

	directorAgent := agents.NewAgent(llm, tools)
	directorAgent.Task("What time is it?")

	ctx := context.TODO()
	executor := agents.NewExecutor(directorAgent, agents.WithShowMessages())
	executor.Run(ctx)

	finalAnswer, _ := directorAgent.GetFinalAnswer()
	fmt.Println(finalAnswer)
}

func (t TimeAgent) Name() string { return "TimeAgent" }

func (t TimeAgent) Call(ctx context.Context, input string) (string, error) {
	currentTime := time.Now()
	fmtCurrentTime := currentTime.Format("2006-01-02 3:04:05 PM")
	return fmtCurrentTime, nil
}
