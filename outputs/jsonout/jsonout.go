// Package json provides a generic JSON output parser for LLM responses.
//
// It unmarshals raw JSON strings into typed Go structs and generates
// instruction prompts that guide an LLM to produce valid JSON output.
package jsonout

import (
	"encoding/json"
	"reflect"
	"strings"
)

// parser is a generic JSON parser that deserializes LLM output into a Go type T.
type parser[T any] struct{}

// NewParser creates a new [parser] for the given type T.
func NewParser[T any]() *parser[T] {
	return &parser[T]{}
}

// Parse deserializes the given JSON string into a value of type T.
// When T is a slice type and the LLM returns a single object, Parse wraps it
// in an array before unmarshaling.
func (p *parser[T]) Parse(output string) (T, error) {
	var result T
	data := strings.TrimSpace(output)
	if reflect.TypeFor[T]().Kind() == reflect.Slice &&
		strings.HasPrefix(data, "{") &&
		strings.HasSuffix(data, "}") {
		data = "[" + data + "]"
	}
	err := json.Unmarshal([]byte(data), &result)
	return result, err
}

// Instructions returns a prompt string that instructs an LLM to produce JSON
// matching the zero-value schema of type T.
func (p *parser[T]) Instructions() string {
	var zero T

	v := reflect.ValueOf(&zero).Elem()
	if v.Kind() == reflect.Slice {
		v.Set(reflect.MakeSlice(v.Type(), 1, 1))
	}

	jsonBytes, err := json.Marshal(zero)
	if err != nil {
		return ""
	}

	return "Output the following JSON schema: " + string(jsonBytes)
}
