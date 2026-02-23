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

	"github.com/bit8bytes/gogantic/agents/tools"
	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
)

type llm interface {
	Generate(ctx context.Context, messages []llms.Message) (*llms.ContentResponse, error)
}

type store interface {
	Add(ctx context.Context, msgs ...llms.Message) error
	List(ctx context.Context) ([]llms.Message, error)
	Clear(ctx context.Context) error
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
	History     store
	actions     []Action
	parser      parser
	finalAnswer string
}

// New creates an agent with the given model, tools, storage, and parser.
// For the ReAct pattern, prefer NewReAct.
func New(model llm, tools []Tool, storage store, p parser) *Agent {
	return &Agent{
		model:   model,
		tools:   toolNames(tools),
		History: storage,
		parser:  p,
	}
}

// Task sets the user's question or task for the agent to solve.
// Call this before starting the Plan-Act loop.
func (a *Agent) Task(ctx context.Context, prompt string) error {
	return a.History.Add(ctx, llms.Message{
		Role:    roles.User,
		Content: "Question: " + prompt,
	})
}

// Plan calls the LLM to decide the next action or provide a final answer.
// Returns Response.Finish=true when the task is complete.
func (a *Agent) Plan(ctx context.Context) (*Response, error) {
	history, err := a.History.List(ctx)
	if err != nil {
		return nil, err
	}

	generated, err := a.model.Generate(ctx, history)
	if err != nil {
		return nil, err
	}

	parsed, err := a.parser.Parse(generated.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent response: %w", err)
	}

	if err := a.addAssistantMessage(ctx, generated.Result); err != nil {
		return nil, fmt.Errorf("failed to store assistant message: %w", err)
	}

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
func (a *Agent) Act(ctx context.Context) error {
	for _, action := range a.actions {
		if err := a.handleAction(ctx, action); err != nil {
			return err
		}
	}
	a.clearActions()
	return nil
}

func (a *Agent) handleAction(ctx context.Context, action Action) error {
	t, exists := a.tools[action.Tool]
	if !exists {
		return a.addObservationMessage(ctx, "The action "+action.Tool+" doesn't exist.")
	}

	observation, err := t.Execute(ctx, tools.Input{
		Content: action.ToolInput,
	})
	if err != nil {
		return a.addObservationMessage(ctx, "Error: "+err.Error())
	}

	return a.addObservationMessage(ctx, observation.Content)
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

func toolNames(tools []Tool) map[string]Tool {
	t := make(map[string]Tool, len(tools))
	for _, tool := range tools {
		t[tool.Name()] = tool
	}
	return t
}
