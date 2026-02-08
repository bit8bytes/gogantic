package agents

type Step struct {
	Thought     string
	Actions     string
	Observation string
}

type Action struct {
	Tool      string
	ToolInput string
	ToolID    string
}

type Finish struct {
	ReturnValues map[string]any
}

type Response struct {
	Actions []Action
	Finish  bool
}
