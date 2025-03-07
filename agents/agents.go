package agents

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bit8bytes/gogantic/core/input"
	"github.com/bit8bytes/gogantic/core/models"
)

type Agent struct {
	Model    models.Model
	Tools    map[string]Tool
	Messages []models.MessageContent
	Actions  []AgentAction
}

func NewAgent(model models.Model, tools map[string]Tool) *Agent {
	toolNames := getToolNames(tools)
	initialMessages := setupReActPromptInitialMessages(toolNames)

	return &Agent{
		Model:    model,
		Tools:    tools,
		Messages: initialMessages,
	}
}

func (a *Agent) Task(prompt string) {
	chatPrompt, _ := input.NewChatPromptTemplate([]models.MessageContent{
		{Role: "user", Content: `
Begin!
Question: {{.input}}
`},
	})

	data := map[string]interface{}{
		"input": prompt,
	}

	formattedMessages, err := chatPrompt.FormatMessages(data)
	if err != nil {
		panic(err)
	}

	a.Messages = append(a.Messages, formattedMessages...)
}

func (a *Agent) Plan(ctx context.Context) (AgentResponse, error) {
	generatedContent, err := a.Model.GenerateContent(ctx, a.Messages)
	if err != nil {
		return AgentResponse{}, err
	}

	text := generatedContent.Result

	if strings.Contains(text, "FINAL ANSWER:") {
		finalAnswerParts := strings.Split(text, "FINAL ANSWER:")
		finalAnswer := strings.TrimSpace(finalAnswerParts[1])

		a.Messages = append(a.Messages, models.MessageContent{
			Role:    "assistant",
			Content: fmt.Sprintf("\nFinal Answer: %s", finalAnswer),
		})

		return AgentResponse{Finish: true}, nil
	}

	reThought := regexp.MustCompile(`(?i)Thought:\s*(.*?)(?:\s*(?:Action:|FINAL ANSWER:|$))`)
	thought := reThought.FindStringSubmatch(text)

	// Extract action - look for format like "Action: [ToolName]"
	reAction := regexp.MustCompile(`Action:\s*\[(.*?)\]`)
	action := reAction.FindStringSubmatch(text)

	// Extract action input - look for format like "Action Input: "input""
	reInput := regexp.MustCompile(`Action Input:\s*"(.*?)"`)
	actionInput := reInput.FindStringSubmatch(text)

	// Handle alternative input format without quotes
	if len(actionInput) <= 1 {
		reInput = regexp.MustCompile(`Action Input:\s*(.*?)(?:\n|$)`)
		actionInput = reInput.FindStringSubmatch(text)
	}

	// Process thought
	if len(thought) > 1 {
		a.addThoughtMessage(strings.TrimSpace(thought[1]))
	}

	// Process action and action input
	if len(action) > 1 {
		a.addActionMessage("[" + action[1] + "]")

		inputText := ""
		if len(actionInput) > 1 {
			inputText = actionInput[1]
			a.addActionInputMessage("\"" + inputText + "\"")
		} else {
			a.addActionInputMessage("\"\"")
		}

		a.Actions = []AgentAction{
			{
				Tool:      action[1],
				ToolInput: inputText,
			},
		}
	} else {
		// If no action is found but we're not at final answer,
		// add a note for debugging
		fmt.Println("Warning: No action found in response")
	}

	return AgentResponse{Finish: false}, nil
}

func (a *Agent) Act(ctx context.Context) {
	for _, action := range a.Actions {
		if !a.handleAction(ctx, action) {
			return
		}
	}
	a.clearActions()
}

func (a *Agent) handleAction(ctx context.Context, action AgentAction) bool {
	tool, exists := a.Tools[action.Tool]
	if !exists {
		a.addObservationMessage("Error: Tool '" + action.Tool + "' not found")
		return false
	}

	observation, err := tool.Call(ctx, action.ToolInput)
	if err != nil {
		a.addObservationMessage("Error: " + err.Error())
		return false
	}

	a.addObservationMessage(observation)
	return true
}

func (a *Agent) clearActions() {
	a.Actions = nil
}

func (a *Agent) GetFinalAnswer() (string, error) {
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

func (a *Agent) CheckIfFinalAnswer(input string) (AgentResponse, error) {
	if strings.Contains(input, "FINAL ANSWER:") {
		finalAnswerParts := strings.Split(input, "FINAL ANSWER:")
		finalAnswer := strings.TrimSpace(finalAnswerParts[1])

		a.Messages = append(a.Messages, models.MessageContent{
			Role:    "assistant",
			Content: fmt.Sprintf("\nFinal Answer: %s", finalAnswer),
		})

		return AgentResponse{Finish: true}, nil
	}

	return AgentResponse{Finish: false}, nil
}

func setupReActPromptInitialMessages(tools string) []models.MessageContent {
	reActPrompt, _ := input.NewChatPromptTemplate([]models.MessageContent{
		{Role: "user", Content: `
Answer the following questions as best you can. 
Use only values from the tools. Do not estimate or predict values.	
Select the tool that fits the question:

[{{.tools}}]

Use the following format:

Thought: you should always think about what to do
Action: the action (only one at a time) to take in suqare braces e.g [NameOfTool]
Action Input: the input value for the action in quotes e.g. "string" or "int"
Observation: the result of the action
... (this Thought:/Action:/Action Input:/Observation: can repeat N times)
Thought: I now know the final answer
FINAL ANSWER: the final answer to the original input question

Think in steps.
`},
	})

	data := map[string]interface{}{
		"tools": tools,
	}

	formattedMessages, err := reActPrompt.FormatMessages(data)
	if err != nil {
		panic(err)
	}

	return formattedMessages
}
