package agents

type AgentStep struct {
	Thought string
	Actions string
	Observation string
}

type AgentAction struct {
	Tool      string
	ToolInput string
	ToolID    string
}

type AgentFinish struct {
	ReturnValues map[string]any
}

type AgentResponse struct {
	Actions []AgentAction
	Finish bool
}