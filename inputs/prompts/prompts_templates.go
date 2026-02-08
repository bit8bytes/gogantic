package prompts

import (
	"bytes"
	"text/template"
)

type prompts struct {
	template *template.Template
}

func New(data string) *prompts {
	tmpl, err := template.New("prompt").Parse(data)
	if err != nil {
		panic(err)
	}
	return &prompts{template: tmpl}
}

func (pt *prompts) Execute(data any) (string, error) {
	var promptBuffer bytes.Buffer
	err := pt.template.Execute(&promptBuffer, data)
	if err != nil {
		return "", err
	}
	return promptBuffer.String(), nil
}
