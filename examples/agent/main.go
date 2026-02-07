package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"github.com/bit8bytes/gogantic/agent"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	tools := []agent.Tool{
		agent.NewWeatherTool(),
		agent.NewTimeTool(),
	}

	agentInstance := agent.New(tools, agent.WithLogger(logger))

	ctx := context.Background()

	log.Println("=== Available Tools ===")
	for _, toolName := range agentInstance.ListTools() {
		description := agentInstance.GetToolDescription(toolName)
		log.Printf("- %s: %s", toolName, description)
	}

	log.Println("\n=== Weather Tool Example ===")
	weatherInput := map[string]interface{}{
		"location": "New York",
		"units":    "metric",
	}
	weatherInputJSON, _ := json.Marshal(weatherInput)

	weatherResult, err := agentInstance.Execute(ctx, "weather", string(weatherInputJSON))
	if err != nil {
		log.Printf("Weather tool error: %v", err)
	} else {
		log.Printf("Weather result: %s", weatherResult)
	}

	log.Println("\n=== Time Tool Example ===")
	timeInput := map[string]interface{}{
		"timezone": "America/New_York",
	}
	timeInputJSON, _ := json.Marshal(timeInput)

	timeResult, err := agentInstance.Execute(ctx, "time", string(timeInputJSON))
	if err != nil {
		log.Printf("Time tool error: %v", err)
	} else {
		log.Printf("Time result: %s", timeResult)
	}

	log.Println("\n=== Error Example (Non-existent Tool) ===")
	_, err = agentInstance.Execute(ctx, "nonexistent", "{}")
	if err != nil {
		log.Printf("Expected error: %v", err)
	}
}
