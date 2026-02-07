# Agent Module

The agent module provides a simple tool-based agent system for gogantic that follows the options pattern for initialization.

## Features

- Tool-based architecture for extensibility
- Options pattern for clean initialization
- Built-in weather and time tools
- Structured logging support
- JSON-based tool input/output

## Usage

```go
import (
    "context"
    "log/slog"
    "github.com/bit8bytes/gogantic/agent"
)

func main() {
    tools := []agent.Tool{
        agent.NewWeatherTool(),
        agent.NewTimeTool(),
    }

    myAgent := agent.New(tools, agent.WithLogger(slog.Default()))

    ctx := context.Background()
    result, err := myAgent.Execute(ctx, "weather", `{"location": "New York"}`)
    if err != nil {
        // handle error
    }
}
```

## Available Tools

### Weather Tool
- **Name**: `weather`
- **Input**: `{"location": string, "units": string}`
- **Description**: Gets weather information for a location

### Time Tool  
- **Name**: `time`
- **Input**: `{"timezone": string}` or `{"location": string}`
- **Description**: Gets current time for a timezone or location

## Creating Custom Tools

```go
type MyTool struct{}

func (t *MyTool) Name() string {
    return "mytool"
}

func (t *MyTool) Description() string {
    return "My custom tool"
}

func (t *MyTool) Execute(ctx context.Context, input string) (string, error) {
    // Tool implementation
    return "result", nil
}
```