package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bit8bytes/gogantic/agent"
	"github.com/bit8bytes/gogantic/llm/openai"
	"github.com/bit8bytes/gogantic/runner"
	"github.com/bit8bytes/gogantic/tool"
)

const (
	ToolGetTemperatureInFahrenheit = "get_temperature_in_fahrenheit"
	ToolFormatFahrenheitToCelsius  = "format_fahrenheit_to_celsius"
)

type GetTemperatureInFahrenheit struct{}
type FormatFahrenheitToCelsius struct{}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("No OPENAI_API_KEY")
	}

	stream := false
	stop := []string{"\nObservation", "Observation"}

	gpt_35_turbo := openai.Model{
		Model:  "gpt-3.5-turbo",
		APIKey: apiKey,
		Stream: &stream,
		Stop:   &stop,
	}

	llm := openai.New(gpt_35_turbo)

	tools := map[string]tool.Tool{
		ToolGetTemperatureInFahrenheit: GetTemperatureInFahrenheit{},
		ToolFormatFahrenheitToCelsius:  FormatFahrenheitToCelsius{},
	}

	weatherAgent := agent.New(llm, tools)
	weatherAgent.Task("What is the temperature outside?")

	runner := runner.New(weatherAgent,
		runner.WithIterationLimit(10),
		runner.WithShowMessages())
	runner.Run(context.TODO())

	finalAnswer, _ := weatherAgent.GetFinalAnswer()
	fmt.Println(finalAnswer)
}

func (t GetTemperatureInFahrenheit) Name() string {
	return ToolGetTemperatureInFahrenheit
}

func (t GetTemperatureInFahrenheit) Schema() string { return `()` }

func (t GetTemperatureInFahrenheit) Call(ctx context.Context, input tool.Input) (tool.Output, error) {
	// This is only for showcase.
	// If you want to use this and handle input e.g. location look at the math agent example.
	return tool.Output{Content: "15.54°F"}, nil
}

func (t FormatFahrenheitToCelsius) Name() string {
	return ToolFormatFahrenheitToCelsius
}

func (t FormatFahrenheitToCelsius) Schema() string { return `(0000)` }

func (t FormatFahrenheitToCelsius) Call(ctx context.Context, input tool.Input) (tool.Output, error) {
	// Still, I do not handle errors in here. This has to be done through testing.
	fahrenheit, _ := strconv.ParseFloat(input.Content, 64)
	celsius := (fahrenheit - 32) * (5.0 / 9.0)
	return tool.Output{Content: fmt.Sprintf("Current temperature: %.2f°C", celsius)}, nil
}
