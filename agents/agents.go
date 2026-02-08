package agents

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bit8bytes/gogantic/agents/tools"
	"github.com/bit8bytes/gogantic/inputs/chats"
	"github.com/bit8bytes/gogantic/llms"
)

type llm interface {
	Generate(ctx context.Context, messages []llms.Message) (*llms.ContentResponse, error)
}

type Tool interface {
	Name() string
	Execute(ctx context.Context, input tools.Input) (tools.Output, error)
}

type Agent struct {
	model    llm
	tools    map[string]Tool
	Messages []llms.Message
	actions  []Action
}

func New(model llm, tools map[string]Tool, messages []llms.Message) *Agent {
	return &Agent{
		model:    model,
		tools:    tools,
		Messages: messages,
	}
}

func (a *Agent) Task(prompt string) error {
	messages := chats.New([]llms.Message{
		{
			Role:    "user",
			Content: "Question: {{.Input}}\n",
		},
	})

	type task struct {
		Input string
	}

	formattedMessages, err := messages.Execute(task{Input: prompt})
	if err != nil {
		return err
	}

	a.Messages = append(a.Messages, formattedMessages...)
	return nil
}

// Identifies the generated messages and splits them into thought, action and action input
func (a *Agent) Plan(ctx context.Context) (*Response, error) {
	generatedContent, err := a.model.Generate(ctx, a.Messages)
	if err != nil {
		return nil, err
	}

	text := generatedContent.Result

	final := extractAfterLabel(text, "FINAL ANSWER:")
	if len(final) > 0 {
		a.Messages = append(a.Messages, llms.Message{
			Role:    "assistant",
			Content: fmt.Sprintf("\nFinal Answer: %s", final),
		})

		return &Response{Finish: true}, nil
	}

	thought := extractAfterLabel(text, "Thought: ")

	// "Action: [ToolName]"
	action := extractAfterLabel(text, "Action: ")

	// "Action Input: "input"
	actionInput := extractAfterLabel(text, "Action Input: ")

	if len(thought) > 1 {
		a.addThoughtMessage(strings.TrimSpace(thought))
	}

	if len(action) > 1 {
		tool := extractSquareBracketsContent(action)
		a.addActionMessage(tool)

		inputText := ""
		if len(actionInput) > 1 {
			inputText = removeQuotes(actionInput)
			a.addActionInputMessage("\"" + inputText + "\"")
		} else {
			a.addActionInputMessage("\"\"")
		}

		a.actions = []Action{
			{
				Tool:      tool,
				ToolInput: inputText,
			},
		}
	} else {
		fmt.Println("Warning: No action found in response")
	}

	return &Response{}, nil
}

// Uses the given tools to get observations
func (a *Agent) Act(ctx context.Context) {
	for _, action := range a.actions {
		if !a.handleAction(ctx, action) {
			return
		}
	}
	a.clearActions()
}

// Handle action is a helper function that calls the tool selected by the LLM and adds the observation output
func (a *Agent) handleAction(ctx context.Context, action Action) bool {
	t, exists := a.tools[action.Tool]
	if !exists {
		a.addObservationMessage("The Action: [" + action.Tool + "] doesn't exist.")
		return false
	}

	i := tools.Input{
		Content: action.ToolInput,
	}

	observation, err := t.Execute(ctx, i)
	if err != nil {
		a.addObservationMessage("Error: " + err.Error())
		return false
	}

	a.addObservationMessage(observation.Content)
	return true
}

func (a *Agent) clearActions() {
	a.actions = nil
}

func (a *Agent) Answer() (string, error) {
	if len(a.Messages) == 0 {
		return "", errors.New("No messages provided")
	}
	finalAnswer := a.Messages[len(a.Messages)-1].Content
	parts := strings.Split(finalAnswer, "Final Answer: ")
	if len(parts) < 2 {
		return "", errors.New("Invalid final answer")
	}
	return parts[1], nil
}

func setupReActPromptInitialMessages(tools string) []llms.Message {
	reActPrompt := chats.New([]llms.Message{
		{Role: "user", Content: `
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
`},
	})

	data := map[string]interface{}{
		"tools": tools,
	}

	formattedMessages, err := reActPrompt.Execute(data)
	if err != nil {
		panic(err)
	}

	return formattedMessages
}

func extractAfterLabel(s, label string) string {
	startIndex := strings.Index(s, label)
	if startIndex == -1 {
		return "" // Label not found
	}
	startIndex += len(label)
	for startIndex < len(s) && s[startIndex] == ' ' {
		startIndex++
	}
	endIndex := strings.Index(s[startIndex:], "\n")
	if endIndex == -1 {
		endIndex = len(s)
	} else {
		endIndex += startIndex
	}

	return s[startIndex:endIndex]
}

func extractSquareBracketsContent(s string) string {
	startIndex := strings.Index(s, "[")
	if startIndex == -1 {
		return "" // No opening bracket found
	}

	endIndex := strings.Index(s[startIndex:], "]")
	if endIndex == -1 {
		return "" // No closing bracket found
	}

	// Extract the content between brackets
	return s[startIndex+1 : startIndex+endIndex]
}

func removeQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
