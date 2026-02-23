package agents

import (
	"context"

	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
)

func (a *Agent) addAssistantMessage(ctx context.Context, content string) error {
	return a.History.Add(ctx, llms.Message{
		Role:    roles.Assistant,
		Content: content,
	})
}

func (a *Agent) addObservationMessage(ctx context.Context, observation string) error {
	return a.History.Add(ctx, llms.Message{
		Role:    roles.System,
		Content: "Observation: " + observation,
	})
}
