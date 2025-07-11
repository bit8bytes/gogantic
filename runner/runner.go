package runner

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/agent"
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

type Runner struct {
	Agent          *agent.Agent
	IterationLimit int
	printMessages  bool
}

type RunnerOption func(*Runner)

func WithIterationLimit(limit int) RunnerOption {
	return func(e *Runner) {
		e.IterationLimit = limit
	}
}

func WithShowMessages() RunnerOption {
	return func(e *Runner) {
		e.printMessages = true
	}
}

func New(agent *agent.Agent, opts ...RunnerOption) *Runner {
	e := &Runner{
		Agent:          agent,
		IterationLimit: 10,
		printMessages:  false,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// Updated Run function for your Runner
func (e *Runner) Run(ctx context.Context) {
	for i := 1; i < e.IterationLimit; i++ {
		todos, err := e.Agent.Plan(ctx)
		if err != nil {
			fmt.Println("Error planning:", err)
			break
		}

		if todos.Finish {
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
