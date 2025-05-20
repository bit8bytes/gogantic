package tool

import (
	"context"
)

type Tool interface {
	Name() string
	Call(ctx context.Context, input string) (string, error)
}
