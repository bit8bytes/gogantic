package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/agents/tools"
	"github.com/bit8bytes/gogantic/inputs/chats"
	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
	"github.com/bit8bytes/gogantic/llms/ollama"
	"github.com/bit8bytes/gogantic/runner"
)

var reActPrompt = `
Answer the following questions as best you can. 
Use only values from the tools. Do not estimate or predict values.	
Select the tool that fits the question:

[{{.tools}}]

Use the following format:

Thought: you should always think about what to do
Action: [Toolname] the action (only one at a time) to take in suqare braces e.g [NameOfTool]
Action Input: "input" the input value for the action in quotes e.g. "value" from Schema
Observation: the result of the action
... (this Thought: .../Action: [Toolname]/Action Input: "input"/Observation: ... can repeat N times)
Thought: I now know the final answer
FINAL ANSWER: the final answer to the original input question

Think in steps. Don't hallucinate. Don't make up answers.
`

type GetTime struct{}

func main() {
	llm := ollama.New(ollama.Model{
		Model:   "gemma3n:e2b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
		Stop:    []string{"\nObservation", "Observation"},
	})

	tools := map[string]agents.Tool{
		"GetTime": GetTime{},
	}

	data := map[string]any{
		"tools": tools,
	}

	messages := []llms.Message{
		{
			Role:    roles.System,
			Content: reActPrompt,
		},
	}
	chats := chats.New(messages)

	formattedMessages, err := chats.Execute(data)
	if err != nil {
		panic(err)
	}

	agent := agents.New(llm, tools, formattedMessages)
	if err := agent.Task("What time is it?"); err != nil {
		panic(err)
	}

	ctx := context.TODO()
	runner := runner.New(agent, runner.WithShowMessages())
	runner.Run(ctx)

	finalAnswer, _ := agent.Answer()
	fmt.Println(finalAnswer)
}

func (t GetTime) Name() string { return "GetTime" }

func (t GetTime) Schema() string { return `()` }

func (t GetTime) Execute(ctx context.Context, input tools.Input) (tools.Output, error) {
	currentTime := time.Now()
	fmtCurrentTime := currentTime.Format("2006-01-02 3:04:05 PM")
	return tools.Output{Content: fmtCurrentTime}, nil
}
