package moderncsqlite

import (
	"context"
	"database/sql"

	"github.com/bit8bytes/gogantic/llms"
	"github.com/bit8bytes/gogantic/stores"
	_ "modernc.org/sqlite"
)

const Sqlite = "sqlite"

func init() {
	stores.Register(Sqlite, NewSqlite)
}

type sqliteStore struct {
	db *sql.DB
}

func NewSqlite(ctx context.Context, name string, db *sql.DB) (stores.Store, error) {
	return &sqliteStore{db: db}, nil
}

func (s *sqliteStore) Add(ctx context.Context, messages ...llms.Message) error {
	for _, msg := range messages {
		_, err := s.db.ExecContext(ctx,
			`INSERT INTO messages (role, content) VALUES (?, ?)`,
			msg.Role.String(), msg.Content,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sqliteStore) List(ctx context.Context) ([]llms.Message, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT role, content FROM messages ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []llms.Message
	for rows.Next() {
		var msg llms.Message
		if err := rows.Scan(&msg.Role, &msg.Content); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, rows.Err()
}

func (s *sqliteStore) Clear(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM messages`)
	return err
}

func (s *sqliteStore) Close() error {
	return s.db.Close()
}
