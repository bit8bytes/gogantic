package agents

import (
	"github.com/bit8bytes/gogantic/inputs/roles"
	"github.com/bit8bytes/gogantic/llms"
)

func (a *Agent) addAssistantMessage(content string) {
	a.Messages = append(a.Messages, llms.Message{
		Role:    roles.Assistent,
		Content: content,
	})
}

func (a *Agent) addObservationMessage(observation string) {
	a.Messages = append(a.Messages, llms.Message{
		Role:    roles.System,
		Content: "Observation: " + observation,
	})
}
