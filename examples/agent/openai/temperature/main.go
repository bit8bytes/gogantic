package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bit8bytes/gogantic/agent"
	"github.com/bit8bytes/gogantic/llm/openai"
)

type CurrentDatetime struct{}
type CurrentTemperatureInFahrenheit struct{}
type FormatFahrenheitToCelsius struct{}
type SaveToFile struct{}

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

	tools := map[string]agent.Tool{
		"CurrentTemperatureInFahrenheit": CurrentTemperatureInFahrenheit{},
		"FormatFahrenheitToCelsius":      FormatFahrenheitToCelsius{},
	}

	weatherAgent := agent.New(llm, tools)
	weatherAgent.Task("What is the temperature outside?")

	executor := agent.NewExecutor(weatherAgent,
		agent.WithIterationLimit(10),
		agent.WithShowMessages())
	executor.Run(context.TODO())

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
