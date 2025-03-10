package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/core/models/llms/ollama"
)

type CurrentDatetime struct{}
type CurrentTemperatureInFahrenheit struct{}
type FormatFahrenheitToCelsius struct{}
type SaveToFile struct{}

func main() {
	mistral_latest := ollama.OllamaModel{
		Model:     "mistral:latest",
		Options:   ollama.ModelOptions{NumCtx: 4096},
		Stream:    false,
		KeepAlive: -1,
		Stop:      []string{"\nObservation", "Observation"}, // Necessary due to the ReAct Prompt Pattern
	}
	llm := ollama.NewOllamaClient(mistral_latest)

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

func (t CurrentTemperatureInFahrenheit) Name() string {
	return "CurrentTemperatureInFahrenheit"
}

func (t CurrentTemperatureInFahrenheit) Call(ctx context.Context, input string) (string, error) {
	// This is only for showcase.
	// If you want to use this and handle input e.g. location look at the math agent example.
	return fmt.Sprintf("5.54°F"), nil
}

func (t FormatFahrenheitToCelsius) Name() string {
	return "FormatFahrenheitToCelsius"
}

func (t FormatFahrenheitToCelsius) Call(ctx context.Context, input string) (string, error) {
	// Still, I do not handle errors in here. This has to be done through testing.
	fahrenheit, _ := strconv.ParseFloat(input, 64)
	celsius := (fahrenheit - 32) * (5.0 / 9.0)
	return fmt.Sprintf("Current temperature: %.2f°C", celsius), nil
}
