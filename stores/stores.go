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

// Memory is a simple in-memory store backed by a slice.
type Memory struct {
	mu   sync.Mutex
	msgs []llms.Message
}

// NewMemory returns a new in-memory store.
func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Add(_ context.Context, msgs ...llms.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.msgs = append(m.msgs, msgs...)
	return nil
}

func (m *Memory) List(_ context.Context) ([]llms.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]llms.Message, len(m.msgs))
	copy(out, m.msgs)
	return out, nil
}

func (m *Memory) Clear(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.msgs = nil
	return nil
}

func (m *Memory) Close() error { return nil }

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
