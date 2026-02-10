// Package agents implements the ReAct (Reasoning and Acting) pattern for LLM-powered agents.
//
// Agents alternate between reasoning about a task and executing tool actions,
// with each observation feeding back into the next reasoning step.
// Use New to create an agent, Task to set the goal, then repeatedly call Plan
// and Act until the agent signals completion.
package agents

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bit8bytes/gogantic/agents/tools"
	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
	"github.com/bit8bytes/gogantic/outputs/jsonout"
)

type llm interface {
	Generate(ctx context.Context, messages []llms.Message) (*llms.ContentResponse, error)
}

// Tool represents an action the agent can perform.
// Each tool must provide a name, description, and execution logic.
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, input tools.Input) (tools.Output, error)
}

type parser interface {
	Parse(text string) (AgentResponse, error)
	Instructions() string
}

// Agent executes tasks using the ReAct pattern (reasoning + acting).
// Call Plan to generate the next action, then Act to execute it.
// Repeat until Plan returns Finish=true, then retrieve the result with Answer.
type Agent struct {
	model       llm
	tools       map[string]Tool
	Messages    []llms.Message
	actions     []Action
	parser      parser
	finalAnswer string
}

// New creates an agent with the given model and tools.
// The agent is initialized with a ReAct system prompt.
func New(model llm, tools []Tool) *Agent {
	p := jsonout.NewParser[AgentResponse]()
	t := toolNames(tools)

	return &Agent{
		model:    model,
		tools:    t,
		Messages: buildReActPrompt(t, p.Instructions()),
		parser:   p,
	}
}

// Task sets the user's question or task for the agent to solve.
// Call this before starting the Plan-Act loop.
func (a *Agent) Task(prompt string) error {
	a.Messages = append(a.Messages, llms.Message{
		Role:    roles.User,
		Content: "Question: " + prompt,
	})
	return nil
}

// Plan calls the LLM to decide the next action or provide a final answer.
// Returns Response.Finish=true when the task is complete.
func (a *Agent) Plan(ctx context.Context) (*Response, error) {
	generated, err := a.model.Generate(ctx, a.Messages)
	if err != nil {
		return nil, err
	}

	parsed, err := a.parser.Parse(generated.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent response: %w", err)
	}

	a.addAssistantMessage(generated.Result)

	if parsed.FinalAnswer != "" {
		a.finalAnswer = parsed.FinalAnswer
		return &Response{Finish: true}, nil
	}

	if parsed.Action == "" {
		return nil, errors.New("agent response contains neither a final answer nor an action")
	}

	a.actions = []Action{
		{
			Tool:      parsed.Action,
			ToolInput: parsed.ActionInput,
		},
	}

	return &Response{}, nil
}

// Act executes the tool chosen by Plan and adds the result as an observation.
// Always call this after Plan (unless Plan returned Finish=true).
func (a *Agent) Act(ctx context.Context) {
	for _, action := range a.actions {
		if !a.handleAction(ctx, action) {
			return
		}
	}
	a.clearActions()
}

func (a *Agent) handleAction(ctx context.Context, action Action) bool {
	t, exists := a.tools[action.Tool]
	if !exists {
		a.addObservationMessage("The action " + action.Tool + " doesn't exist.")
		return false
	}

	observation, err := t.Execute(ctx, tools.Input{
		Content: action.ToolInput,
	})
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

// Answer returns the final result after the agent completes the task.
// Only call this after Plan returns Finish=true.
func (a *Agent) Answer() (string, error) {
	if a.finalAnswer == "" {
		return "", errors.New("no final answer available")
	}
	return a.finalAnswer, nil
}

func buildReActPrompt(tools map[string]Tool, jsonInstructions string) []llms.Message {
	var toolDescriptions strings.Builder
	for _, t := range tools {
		fmt.Fprintf(&toolDescriptions, "- %s: %s\n", t.Name(), t.Description())
	}

	return []llms.Message{
		{
			Role: roles.System,
			Content: fmt.Sprintf(`
You are an helpful agent. Answer questions using the available tools.
Do not estimate or predict values. Use only values returned by tools.

Available tools:
%s
%s

Respond with a JSON object on each turn with these fields:
- "thought": your reasoning about what to do next
- "action": the exact tool name to call (empty string when giving final answer)
- "action_input": the input to pass to the tool (empty string when giving final answer)
- "final_answer": your final answer (empty string when calling a tool)

Think step by step. Do not hallucinate.`, toolDescriptions.String(), jsonInstructions),
		},
	}
}

func toolNames(tools []Tool) map[string]Tool {
	t := make(map[string]Tool, len(tools))
	for _, tool := range tools {
		t[tool.Name()] = tool
	}
	return t
}
