package json

import (
	"encoding/json"
	"reflect"

	"github.com/bit8bytes/gogantic/input/prompt"
)

type Parser[T any] struct{}

func (p *Parser[T]) Parse(output string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(output), &result)
	return result, err
}

func (p *Parser[T]) ParseWithPrompt(output string, prompt prompt.Prompt) (T, error) {
	// Implement logic if necessary, otherwise just parse the output
	return p.Parse(output)
}

func (p *Parser[T]) GetFormatInstructions() string {
	// Get the JSON schema as a string
	jsonSchema := getJSONSchema(reflect.TypeOf((*T)(nil)).Elem())

	return "Return the output as JSON with schema: " + jsonSchema
}

func getJSONSchema(t reflect.Type) string {
	if t.Kind() != reflect.Struct {
		return ""
	}

	fields := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			fields[field.Name] = jsonTag
		}
	}

	jsonStr, err := json.Marshal(fields)
	if err != nil {
		return ""
	}

	return string(jsonStr)
}
