// Package runner provides execution orchestration for agentic Plan-Act loops.
// It manages the iterative cycle of planning and acting until task completion
// or iteration limits are reached.
package runner

import (
	"context"
	"fmt"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/inputs/roles"
)

const (
	reset = "\033[0m"
	blue  = "\033[34m"
	green = "\033[32m"
	cyan  = "\033[36m"
	white = "\033[37m"
)

// runner executes an agent's Plan-Act loop.
type runner struct {
	agent             *agents.Agent
	printMessages     bool
	lastPrintedMsgIdx int
}

// New creates a Runner with the given agent, iteration limit, and message printing option.
func New(agent *agents.Agent, printMessages bool) *runner {
	return &runner{
		agent:             agent,
		printMessages:     printMessages,
		lastPrintedMsgIdx: 0,
	}
}

// Run executes the agent's Plan-Act loop until completion or iteration limit.
func (r *runner) Run(ctx context.Context) error {
RUN:
	for {
		select {
		case <-ctx.Done():
			break RUN
		default:
			response, err := r.agent.Plan(ctx)
			if err != nil {
				return fmt.Errorf("planning failed: %w", err)
			}

			if response.Finish {
				return nil
			}

			r.agent.Act(ctx)

			r.printNewMessages()
		}
	}
	return fmt.Errorf("no final answer available")
}

func (r *runner) printNewMessages() {
	if !r.printMessages {
		return
	}

	messages := r.agent.Messages
	startIdx := r.lastPrintedMsgIdx + 1

	for i := startIdx; i < len(messages); i++ {
		msg := messages[i]
		color := r.getColorForRole(msg.Role)
		fmt.Printf("%s%s: %s%s\n", color, msg.Role, msg.Content, reset)
	}

	r.lastPrintedMsgIdx = len(messages) - 1
}

func (r *runner) getColorForRole(role roles.Role) string {
	switch role {
	case roles.Assistent:
		return blue
	case roles.System:
		return green
	case roles.User:
		return cyan
	default:
		return white
	}
}
