// Package json provides a generic JSON output parser for LLM responses.
//
// It unmarshals raw JSON strings into typed Go structs and generates
// instruction prompts that guide an LLM to produce valid JSON output.
package json

import (
	"encoding/json"
)

const (
	Format string = "json"
)

// parser is a generic JSON parser that deserializes LLM output into a Go type T.
type parser[T any] struct{}

// NewParser creates a new [parser] for the given type T.
func NewParser[T any]() *parser[T] {
	return &parser[T]{}
}

// Parse deserializes the given JSON string into a value of type T.
func (p *parser[T]) Parse(output string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(output), &result)
	return result, err
}

// Instructions returns a prompt string that instructs an LLM to produce JSON
// matching the zero-value schema of type T.
func (p *parser[T]) Instructions() string {
	var zero T
	jsonBytes, err := json.Marshal(zero)
	if err != nil {
		return ""
	}
	return "Output the following JSON schema: " + string(jsonBytes)
}
