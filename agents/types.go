package agents

// AgentResponse is the JSON schema the LLM produces each iteration.
type AgentResponse struct {
	Thought     string `json:"thought"`
	Action      string `json:"action"`
	ActionInput string `json:"action_input"`
	FinalAnswer string `json:"final_answer"`
}

// Action is the internal representation of a tool call extracted from AgentResponse.
type Action struct {
	Tool      string
	ToolInput string
}

// Response indicates whether the agent loop should continue or finish.
type Response struct {
	Actions []Action
	Finish  bool
}
