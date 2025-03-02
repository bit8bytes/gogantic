package input

import (
    "testing"
)

func TestNewPromptTemplate(t *testing.T) {
    validTemplateString := "Hello, {{.name}}!"
    invalidTemplateString := "Hello, {{.name}"

    tests := []struct {
        name        string
        templateStr string
        shouldError bool
    }{
        {"Valid Template", validTemplateString, false},
        {"Invalid Template", invalidTemplateString, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            pt, err := NewPromptTemplate(tt.templateStr)
            if tt.shouldError {
                if err == nil {
                    t.Errorf("expected error, got nil")
                }
                return
            }
            if err != nil {
                t.Errorf("did not expect error, got %v", err)
            }
            if pt.Template == nil {
                t.Errorf("expected template to be initialized, got nil")
            }
        })
    }
}

func TestFormat(t *testing.T) {
    tmplStr := "Hello, {{.Name}}!"
    pt, err := NewPromptTemplate(tmplStr)
    if err != nil {
        t.Fatalf("unexpected error creating template: %v", err)
    }

    tests := []struct {
        name       string
        data       interface{}
        want       string
        shouldFail bool
    }{
        {"Valid Data", map[string]string{"Name": "World"}, "Hello, World!", false},
        {"Missing Field", map[string]string{"WrongField": "World"}, "Hello, <no value>!", false},
        {"Nil Data", nil, "Hello, <no value>!", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := pt.Format(tt.data)
            if tt.shouldFail {
                if err == nil {
                    t.Errorf("expected error, got nil")
                }
                return
            }
            if err != nil {
                t.Errorf("did not expect error, got %v", err)
            }
            if got != tt.want {
                t.Errorf("expected %q, got %q", tt.want, got)
            }
        })
    }
}