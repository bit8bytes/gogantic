package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bit8bytes/gogantic/agent"
	"github.com/bit8bytes/gogantic/llm/ollama"
	"github.com/bit8bytes/gogantic/runner"
	"github.com/bit8bytes/gogantic/tool"
)

type CurrentDatetime struct{}
type CurrentTemperatureInFahrenheit struct{}
type FormatFahrenheitToCelsius struct{}
type SaveToFile struct{}

func main() {
	mistral_latest := ollama.Model{
		Model:     "gemma3:4b",
		Options:   ollama.Options{NumCtx: 4096},
		Stream:    false,
		KeepAlive: -1,
		Stop:      []string{"\nObservation", "Observation"}, // Necessary due to the ReAct Prompt Pattern
	}
	llm := ollama.New(mistral_latest)

	tools := map[string]tool.Tool{
		"CurrentTemperatureInFahrenheit": CurrentTemperatureInFahrenheit{},
		"FormatFahrenheitToCelsius":      FormatFahrenheitToCelsius{},
	}

	weatherAgent := agent.New(llm, tools)
	weatherAgent.Task("1. What is the temperature outside? 2. What is the temperature in Celsius?")

	runner := runner.New(weatherAgent,
		runner.WithIterationLimit(10),
		runner.WithShowMessages())
	runner.Run(context.TODO())

	finalAnswer, _ := weatherAgent.GetFinalAnswer()
	fmt.Println(finalAnswer)
}

func (t CurrentTemperatureInFahrenheit) Name() string {
	return "CurrentTemperatureInFahrenheit"
}

func (t CurrentTemperatureInFahrenheit) Call(ctx context.Context, input string) (string, error) {
	// This is only for showcase.
	// If you want to use this and handle input e.g. location look at the math agent example.
	return fmt.Sprintf("15.54°F"), nil
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
