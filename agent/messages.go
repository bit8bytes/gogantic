package agent

import (
	"fmt"

	"github.com/bit8bytes/gogantic/llm"
)

func (a *Agent) addObservationMessage(observation string) {
	a.Messages = append(a.Messages, llm.Message{
		Role:    "system", // Use system role for observations
		Content: "Observation: " + observation,
	})
}

// Helper method to add thought message
func (a *Agent) addThoughtMessage(thought string) {
	a.Messages = append(a.Messages, llm.Message{
		Role:    "assistant",
		Content: "Thought: " + thought,
	})
}

// Helper method to add action message
func (a *Agent) addActionMessage(action string) {
	a.Messages = append(a.Messages, llm.Message{
		Role:    "assistant",
		Content: fmt.Sprintf(`Action: [%s]`, action),
	})
}

// Helper method to add action input message
func (a *Agent) addActionInputMessage(input string) {
	a.Messages = append(a.Messages, llm.Message{
		Role:    "assistant",
		Content: `Action Input: ` + input,
	})
}
