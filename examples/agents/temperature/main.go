package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/core/models/llms/ollama"
)

// Implemented tools at the bottom of this file
type CurrentDatetime struct{}
type CurrentTemperatureInFahrenheit struct{}
type FormatFahrenheitToCelsius struct{}
type SaveToFile struct{}

func main() {
	// Agent need a llm and tools
	// Agent Executor iterates 10 times that the agent can solve the task.
	// Currently the agent works with the ReAct Prompt Pattern
	phi3 := ollama.OllamaModel{
		Model:     "mistral:latest", // This is the best working model from ollama, currently.
		Options:   ollama.ModelOptions{NumCtx: 4096},
		Stream:    false,
		KeepAlive: -1,
		Stop:      []string{"\nObservation", "Observation"}, // Necessary due to the ReAct Prompt Pattern
	}
	llm := ollama.NewOllamaClient(phi3)

	tools := map[string]agents.Tool{
		"CurrentTemperatureInFahrenheit": CurrentTemperatureInFahrenheit{},
		"FormatFahrenheitToCelsius":      FormatFahrenheitToCelsius{},
	}

	weatherAgent := agents.NewAgent(llm, tools)
	weatherAgent.Task("What is the temperature outside?")

	ctx := context.TODO()

	executor := agents.NewExecutor(weatherAgent,
		agents.WithIterationLimit(20),
		agents.WithShowMessages())
	executor.Run(ctx)

	finalAnswer, _ := weatherAgent.GetFinalAnswer()
	fmt.Println(finalAnswer)
}

// Implementation of the tools!
// Tools follow the tools interface with Name and Call.
func (t CurrentTemperatureInFahrenheit) Name() string {
	return "CurrentTemperatureInFahrenheit"
}

func (t CurrentTemperatureInFahrenheit) Call(ctx context.Context, input string) (string, error) {
	return fmt.Sprintf("5.54°F"), nil
}

func (t FormatFahrenheitToCelsius) Name() string {
	return "FormatFahrenheitToCelsius"
}

func (t FormatFahrenheitToCelsius) Call(ctx context.Context, input string) (string, error) {
	fahrenheit, _ := strconv.ParseFloat(input, 64)
	celsius := (fahrenheit - 32) * (5.0 / 9.0)
	return fmt.Sprintf("Current temperature: %.2f°C", celsius), nil
}
