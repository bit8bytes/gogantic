// Package prompts provides templated prompt handling for LLM interactions.
package prompts

import (
	"bytes"
	"text/template"
)

type prompts struct {
	template *template.Template
}

// New creates a new prompts instance by parsing the given string as a
// Go text/template. It panics if the string fails to parse.
func New(data string) *prompts {
	tmpl, err := template.New("prompt").Parse(data)
	if err != nil {
		panic(err)
	}
	return &prompts{template: tmpl}
}

// Execute applies the given data to the template and returns the resulting string.
func (pt *prompts) Execute(data any) (string, error) {
	var promptBuffer bytes.Buffer
	err := pt.template.Execute(&promptBuffer, data)
	if err != nil {
		return "", err
	}
	return promptBuffer.String(), nil
}
