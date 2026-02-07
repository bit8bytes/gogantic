package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type WeatherTool struct{}

func NewWeatherTool() *WeatherTool {
	return &WeatherTool{}
}

func (w *WeatherTool) Name() string {
	return "weather"
}

func (w *WeatherTool) Description() string {
	return "Get weather information for a location"
}

func (w *WeatherTool) Execute(ctx context.Context, input string) (string, error) {
	var params struct {
		Location string `json:"location"`
		Units    string `json:"units,omitempty"`
	}

	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", fmt.Errorf("%w: %v", ErrToolInput, err)
	}

	if params.Location == "" {
		return "", fmt.Errorf("%w: location is required", ErrToolInput)
	}

	if params.Units == "" {
		params.Units = "metric"
	}

	result := map[string]interface{}{
		"location": params.Location,
		"temperature": map[string]interface{}{
			"value": 22,
			"units": params.Units,
			"scale": "Celsius",
		},
		"condition": "Partly Cloudy",
		"humidity":  65,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	response, _ := json.Marshal(result)
	return string(response), nil
}

type TimeTool struct{}

func NewTimeTool() *TimeTool {
	return &TimeTool{}
}

func (t *TimeTool) Name() string {
	return "time"
}

func (t *TimeTool) Description() string {
	return "Get current time for a timezone or location"
}

func (t *TimeTool) Execute(ctx context.Context, input string) (string, error) {
	var params struct {
		Timezone string `json:"timezone,omitempty"`
		Location string `json:"location,omitempty"`
	}

	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", fmt.Errorf("%w: %v", ErrToolInput, err)
	}

	var location string
	if params.Timezone != "" {
		location = params.Timezone
	} else if params.Location != "" {
		location = params.Location
	} else {
		location = "UTC"
	}

	var timeStr string
	switch location {
	case "UTC":
		timeStr = time.Now().UTC().Format(time.RFC3339)
	case "America/New_York", "New York":
		loc, _ := time.LoadLocation("America/New_York")
		timeStr = time.Now().In(loc).Format(time.RFC3339)
	case "Europe/London", "London":
		loc, _ := time.LoadLocation("Europe/London")
		timeStr = time.Now().In(loc).Format(time.RFC3339)
	case "Asia/Tokyo", "Tokyo":
		loc, _ := time.LoadLocation("Asia/Tokyo")
		timeStr = time.Now().In(loc).Format(time.RFC3339)
	default:
		timeStr = time.Now().UTC().Format(time.RFC3339)
	}

	result := map[string]interface{}{
		"time":      timeStr,
		"timezone":  location,
		"timestamp": time.Now().Unix(),
	}

	response, _ := json.Marshal(result)
	return string(response), nil
}
