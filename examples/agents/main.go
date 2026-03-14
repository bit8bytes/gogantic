package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bit8bytes/beago/agents"
	"github.com/bit8bytes/beago/llms/ollama"
	"github.com/bit8bytes/beago/runner"
	"github.com/bit8bytes/beago/stores/moderncsqlite"
	"github.com/bit8bytes/beago/tools"
)

func setupDatabase(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "agent.db")
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDatabase(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS messages (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			role       TEXT NOT NULL,
			content    TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`)
	return err
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	// Omitting errors for example sake.
	db, _ := setupDatabase(ctx)
	defer db.Close()

	_ = migrateDatabase(ctx, db)

	storage, _ := moderncsqlite.New(ctx, db)
	defer storage.Close()

	model := ollama.New(ollama.Model{
		Model:   "gemma3:12b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
		Format:  ollama.JSON,
	})

	// These tools are specifically designed for Golang.
	tools := []agents.Tool{
		tools.RunGoVet{},
	}

	agent, err := agents.NewReAct(ctx, model, tools, storage)
	if err != nil {
		panic(err)
	}

	task := `Where is the NewReAct function called?`
	if err := agent.Task(ctx, task); err != nil {
		panic(err)
	}

	r := runner.New(agent, true)
	if err := r.Run(ctx); err != nil {
		panic(err)
	}

	finalAnswer, err := agent.Answer()
	if errors.Is(err, agents.ErrNoFinalAnswer) {
		fmt.Println("No final answer found")
		return
	}
	fmt.Println(finalAnswer)
}
