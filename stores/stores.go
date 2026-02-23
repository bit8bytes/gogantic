package stores

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

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

// Opener defines the functional factory for creating a Store backed by a sql.DB.
type Opener func(ctx context.Context, name string, db *sql.DB) (Store, error)

var (
	mu      sync.RWMutex
	drivers = make(map[string]Opener)
)

// Register makes a store driver available by the provided name.
// If Register is called twice with the same name or if driver is nil, it panics.
func Register(name string, fn Opener) {
	mu.Lock()
	defer mu.Unlock()
	if fn == nil {
		panic("store: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("store: Register called twice for driver " + name)
	}
	drivers[name] = fn
}

// New opens a store specified by its driver name.
func New(ctx context.Context, name string, db *sql.DB) (Store, error) {
	mu.RLock()
	fn, ok := drivers[name]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("store: unknown driver %q (forgotten import?)", name)
	}
	return fn(ctx, name, db)
}
