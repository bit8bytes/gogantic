package stores

import (
	"context"

	"github.com/bit8bytes/gogantic/llms"
)

// Store is the interface that wraps the basic message persistence methods.
//
// Implementations of Store must be safe for concurrent use by multiple goroutines.
type Store interface {
	Add(ctx context.Context, msgs ...llms.Message) error
	List(ctx context.Context) ([]llms.Message, error)
	Clear(ctx context.Context) error
	Close() error
}
