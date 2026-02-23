package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
	"github.com/bit8bytes/gogantic/outputs/jsonout"
)

// NewReAct creates an agent pre-configured for the ReAct pattern.
// It seeds the ReAct system prompt into storage.
func NewReAct(ctx context.Context, model llm, tools []Tool, storage store) (*Agent, error) {
	p := jsonout.NewParser[AgentResponse]()
	t := toolNames(tools)

	msgs := buildReActPrompt(t, p.Instructions())
	if err := storage.Add(ctx, msgs...); err != nil {
		return nil, err
	}

	return &Agent{
		model:   model,
		tools:   t,
		History: storage,
		parser:  p,
	}, nil
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
