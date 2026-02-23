package agents

import (
	"context"

	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
)

func (a *Agent) addAssistantMessage(ctx context.Context, content string) {
	a.History.Add(ctx, llms.Message{
		Role:    roles.Assistent,
		Content: content,
	})
}

func (a *Agent) addObservationMessage(ctx context.Context, observation string) {
	a.History.Add(ctx, llms.Message{
		Role:    roles.System,
		Content: "Observation: " + observation,
	})
}
