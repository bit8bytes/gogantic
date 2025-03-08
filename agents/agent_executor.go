package agents

import (
	"context"
	"fmt"
)

type Executor struct {
	Agent          *Agent
	IterationLimit int
	printMessages  bool
}

type ExecutorOption func(*Executor)

func WithIterationLimit(limit int) ExecutorOption {
	return func(e *Executor) {
		e.IterationLimit = limit
	}
}

func WithShowMessages() ExecutorOption {
	return func(e *Executor) {
		e.printMessages = true
	}
}

func NewExecutor(agent *Agent, opts ...ExecutorOption) *Executor {
	e := &Executor{
		Agent:          agent,
		IterationLimit: 10,
		printMessages:  false,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// Updated Run function for your Executor
func (e *Executor) Run(ctx context.Context) {
	for i := 1; i < e.IterationLimit; i++ {
		todos, err := e.Agent.Plan(ctx)
		if err != nil {
			fmt.Println("Error planning:", err)
			break
		}

		if todos.Finish {
			fmt.Println("Found final answer!")
			break
		}

		e.Agent.Act(ctx)

		if e.printMessages && len(e.Agent.Messages) > 0 {
			thought := fmt.Sprintf("%s: %s", e.Agent.Messages[len(e.Agent.Messages)-4].Role, e.Agent.Messages[len(e.Agent.Messages)-4].Content)
			action := fmt.Sprintf("%s: %s", e.Agent.Messages[len(e.Agent.Messages)-3].Role, e.Agent.Messages[len(e.Agent.Messages)-3].Content)
			actionInput := fmt.Sprintf("%s: %s", e.Agent.Messages[len(e.Agent.Messages)-2].Role, e.Agent.Messages[len(e.Agent.Messages)-2].Content)
			observation := fmt.Sprintf("%s: %s", e.Agent.Messages[len(e.Agent.Messages)-1].Role, e.Agent.Messages[len(e.Agent.Messages)-1].Content)
			fmt.Println(Blue + thought + Reset)
			fmt.Println(Yellow + action + Reset)
			fmt.Println(Yellow + actionInput + Reset)
			fmt.Println(Green + observation + Reset)
		}
	}
}
