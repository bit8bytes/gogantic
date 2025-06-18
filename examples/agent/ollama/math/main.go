package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/bit8bytes/gogantic/agent"
	"github.com/bit8bytes/gogantic/llm/ollama"
	"github.com/bit8bytes/gogantic/output"
	"github.com/bit8bytes/gogantic/output/separator"
	"github.com/bit8bytes/gogantic/runner"
	"github.com/bit8bytes/gogantic/tool"
)

type Calculator struct{}

func main() {
	mistral_latest := ollama.Model{
		Model:   "mistral:latest",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
		Stop:    []string{"\nObservation", "Observation"},
	}
	llm := ollama.New(mistral_latest)

	tools := map[string]tool.Tool{
		"Calculator": Calculator{},
	}

	// The math agent does wild calculator calls. But at the end somehow it comes to an result.
	// It depends on how the tool instrcuts the model to handle the input
	mathAgent := agent.New(llm, tools)
	mathAgent.Task("What is 22 * 13?")

	ctx := context.TODO()
	runner := runner.New(mathAgent, runner.WithShowMessages())
	runner.Run(ctx)

	finalAnswer, _ := mathAgent.GetFinalAnswer()
	fmt.Println(finalAnswer)
}

func (c Calculator) Name() string { return "Calculator" }

func (c Calculator) Call(ctx context.Context, input string) (string, error) {
	var seperator output.Parser[[]string] = &separator.Space{}
	seperatedResult, _ := seperator.Parse(input)

	if len(seperatedResult) == 3 {
		operation := seperatedResult[0]
		a, _ := strconv.ParseFloat(seperatedResult[1], 64)
		b, _ := strconv.ParseFloat(seperatedResult[2], 64)

		var result float64
		switch operation {
		case "add":
			result = a + b
		case "sub":
			result = a - b
		case "mul":
			result = a * b
		case "div":
			if b == 0 {
				return "", errors.New("Cannot divide by zero")
			}
			result = a / b
		case "mod":
			result = float64(int(a) % int(b))
		default:
			return "", errors.New("Invalid operation. Supported operations are: add, sub, mul, div, mod")
		}

		return fmt.Sprintf("%.2f", result), nil
	}

	// The error message could be a place to evaluate by another agent what went wrong.
	return "", errors.New("Please provide space separted input e.g. \"add a b\" or \"sub a b\"")
}
