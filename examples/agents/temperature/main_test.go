package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/core/models/llms/ollama"
)

// TestTemperatureConversions tests different temperature conversions
func TestTemperatureConversions(t *testing.T) {
	tests := []struct {
		name          string
		fahrenheit    float64
		expectedRange [2]float64
	}{
		{"Freezing", 32, [2]float64{-0.1, 0.1}},            // 32°F = 0°C
		{"Room Temperature", 72, [2]float64{21.9, 22.3}},   // 72°F = 22.22°C
		{"Body Temperature", 98.6, [2]float64{36.9, 37.1}}, // 98.6°F = 37°C
		{"Hot Day", 100, [2]float64{37.7, 37.9}},           // 100°F = 37.78°C
		{"Very Cold", -40, [2]float64{-40.1, -39.9}},       // -40°F = -40°C
		{"Very Hot", 212, [2]float64{99.9, 100.1}},         // 212°F = 100°C
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempTool := MockCurrentTemperatureInFahrenheit{temp: tc.fahrenheit}
			conversionTool := FormatFahrenheitToCelsius{}

			phi3 := ollama.OllamaModel{
				Model:     "mistral:latest",
				Options:   ollama.ModelOptions{NumCtx: 4096},
				Stream:    false,
				KeepAlive: -1,
				Stop:      []string{"\nObservation", "Observation"},
			}
			llm := ollama.NewOllamaClient(phi3)

			tools := map[string]agents.Tool{
				"CurrentTemperatureInFahrenheit": tempTool,
				"FormatFahrenheitToCelsius":      conversionTool,
			}

			weatherAgent := agents.NewAgent(llm, tools)
			weatherAgent.Task("What is the temperature outside?")

			ctx := context.TODO()
			executor := agents.NewExecutor(weatherAgent)
			executor.Run(ctx)

			finalAnswer, err := weatherAgent.GetFinalAnswer()
			if err != nil {
				t.Fatalf("Failed to get final answer: %v", err)
			}

			fmt.Printf("Test case: %s, Final answer: %s\n", tc.name, finalAnswer)

			// Extract temperature from the final answer using regex
			celsiusValue, found := extractTemperature(finalAnswer)
			if !found {
				t.Errorf("Could not find a temperature value in the answer: %s", finalAnswer)
			} else {
				// Check if extracted temperature is within the expected range
				if celsiusValue < tc.expectedRange[0] || celsiusValue > tc.expectedRange[1] {
					t.Errorf("Extracted temperature %.2f°C not within expected range [%.2f, %.2f]",
						celsiusValue, tc.expectedRange[0], tc.expectedRange[1])
				} else {
					fmt.Printf("Successfully verified temperature %.2f°C is within range [%.2f, %.2f]\n",
						celsiusValue, tc.expectedRange[0], tc.expectedRange[1])
				}
			}

			// Also test the tool directly to verify conversion accuracy
			directResult, err := conversionTool.Call(ctx, fmt.Sprintf("%.1f", tc.fahrenheit))
			if err != nil {
				t.Fatalf("Direct tool call failed: %v", err)
			}

			// Extract the celsius value from the direct tool result
			var celsius float64
			fmt.Sscanf(directResult, "Current temperature: %f°C", &celsius)

			// Check if it's within the expected range
			if celsius < tc.expectedRange[0] || celsius > tc.expectedRange[1] {
				t.Errorf("Direct conversion result %.2f°C not within expected range [%.2f, %.2f]",
					celsius, tc.expectedRange[0], tc.expectedRange[1])
			}
		})
	}
}

// extractTemperature extracts a temperature value from a string
// It looks for patterns like "X°C", "X degrees Celsius", etc.
func extractTemperature(s string) (float64, bool) {
	// Look for patterns like X.XX°C or X.XX degrees Celsius
	re := regexp.MustCompile(`(-?\d+\.?\d*)(?:\s*(?:°C|degrees Celsius|degrees C))`)
	matches := re.FindStringSubmatch(s)

	if len(matches) > 1 {
		temp, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return temp, true
		}
	}

	// Try another pattern without the degree symbol
	re = regexp.MustCompile(`(-?\d+\.?\d*)(?:\s*(?:C|celsius))`)
	matches = re.FindStringSubmatch(strings.ToLower(s))

	if len(matches) > 1 {
		temp, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return temp, true
		}
	}

	// As a fallback, try to extract any number
	re = regexp.MustCompile(`(-?\d+\.?\d+)`)
	matches = re.FindStringSubmatch(s)

	if len(matches) > 1 {
		temp, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return temp, true
		}
	}

	return 0, false
}

// MockCurrentTemperatureInFahrenheit returns a specific temperature
type MockCurrentTemperatureInFahrenheit struct {
	temp float64
}

func (t MockCurrentTemperatureInFahrenheit) Name() string {
	return "CurrentTemperatureInFahrenheit"
}

func (t MockCurrentTemperatureInFahrenheit) Call(ctx context.Context, input string) (string, error) {
	return fmt.Sprintf("%.1f°F", t.temp), nil
}

// cleanTemperatureString removes °F or other non-numeric characters
func cleanTemperatureString(s string) string {
	var result string
	for i, c := range s {
		if (c >= '0' && c <= '9') || c == '.' || (i == 0 && c == '-') {
			result += string(c)
		}
	}
	return result
}
