package input

import (
    "bytes"
    "text/template"
)

type PromptTemplate struct {
	Template *template.Template
}

func NewPromptTemplate(templateString string) (*PromptTemplate, error) {
    tmpl, err := template.New("go-template").Parse(templateString)
    if err != nil {
        return nil, err
    }
    return &PromptTemplate{Template: tmpl}, nil
}

func (pt *PromptTemplate) Format(data interface{}) (string, error) {
    var promptBuffer bytes.Buffer
    err := pt.Template.Execute(&promptBuffer, data)
    if err != nil {
        return "", err
    }
    return promptBuffer.String(), nil
}