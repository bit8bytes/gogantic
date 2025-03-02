package output

import (
	"encoding/json"
	"reflect"
)

type JsonOutputParser[T any] struct{}

func (p *JsonOutputParser[T]) Parse(output string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(output), &result)
	return result, err
}

func (p *JsonOutputParser[T]) ParseWithPrompt(output string, prompt PromptValue) (T, error) {
	// Implement logic if necessary, otherwise just parse the output
	return p.Parse(output)
}

func (p *JsonOutputParser[T]) GetFormatInstructions() string {
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
