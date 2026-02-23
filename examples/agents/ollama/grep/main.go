package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/llms/ollama"
	"github.com/bit8bytes/gogantic/runner"
	"github.com/bit8bytes/gogantic/stores"
	modernc "github.com/bit8bytes/gogantic/stores/moderncsqlite"
	_ "github.com/bit8bytes/gogantic/stores/moderncsqlite"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	db, err := sql.Open(modernc.Sqlite, ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS messages (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			role       TEXT NOT NULL,
			content    TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`)
	if err != nil {
		panic(err)
	}

	storage, err := stores.New(ctx, modernc.Sqlite, db)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	model := ollama.New(ollama.Model{
		Model:   "gemma3n:e2b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
		Format:  "json",
	})

	agent, err := agents.NewReAct(ctx, model, []agents.Tool{ListDir{}}, storage)
	if err != nil {
		panic(err)
	}

	if err := agent.Task(ctx, "List all files in folder agents/"); err != nil {
		panic(err)
	}

	r := runner.New(agent, true)
	if err := r.Run(ctx); err != nil {
		panic(err)
	}

	finalAnswer, err := agent.Answer()
	if err != nil {
		panic(err)
	}
	fmt.Println(finalAnswer)
}
