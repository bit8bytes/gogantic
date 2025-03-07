package agents

import (
	"fmt"

	"github.com/bit8bytes/gogantic/core/models"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
)

func (a *Agent) addObservationMessage(observation string) {
	a.Messages = append(a.Messages, models.MessageContent{
		Role:    "system", // Use system role for observations
		Content: "Observation: " + observation,
	})
}

// Helper method to add thought message
func (a *Agent) addThoughtMessage(thought string) {
	a.Messages = append(a.Messages, models.MessageContent{
		Role:    "assistant",
		Content: "Thought: " + thought,
	})
}

// Helper method to add action message
func (a *Agent) addActionMessage(action string) {
	a.Messages = append(a.Messages, models.MessageContent{
		Role:    "assistant",
		Content: "Action: " + action,
	})
}

// Helper method to add action input message
func (a *Agent) addActionInputMessage(input string) {
	a.Messages = append(a.Messages, models.MessageContent{
		Role:    "assistant",
		Content: "Action Input: " + input,
	})
}

func (a *Agent) addObservationError(err error) {
	message := models.MessageContent{
		Role:    "assistant",
		Content: fmt.Sprintf("Observation: %s", err),
	}
	a.Messages = append(a.Messages, message)
}

func (a *Agent) addMissingTool() {
	tools := getToolNames(a.Tools)
	message := models.MessageContent{
		Role:    "user",
		Content: fmt.Sprintf("Please use one of these tools [%s]\n", tools),
	}
	a.Messages = append(a.Messages, message)
}
