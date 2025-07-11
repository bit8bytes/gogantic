package tool

import (
	"context"
)

type Input struct {
	Content string
}

type Output struct {
	Content string
}

type Tool interface {
	Name() string
	Call(ctx context.Context, input Input) (Output, error)
	Schema() string
}
