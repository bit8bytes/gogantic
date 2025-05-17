package prompt

import (
	"bytes"
	"text/template"
)

type Template struct {
	Template *template.Template
}

type Prompt struct {
	Content string
}

func New(templateString string) (*Template, error) {
	tmpl, err := template.New("go-template").Parse(templateString)
	if err != nil {
		return nil, err
	}
	return &Template{Template: tmpl}, nil
}

func (pt *Template) Format(data interface{}) (string, error) {
	var promptBuffer bytes.Buffer
	err := pt.Template.Execute(&promptBuffer, data)
	if err != nil {
		return "", err
	}
	return promptBuffer.String(), nil
}
