package prompts

import (
	"testing"
)

func TestNewPromptTemplate(t *testing.T) {
	validTemplateString := "Hello, {{.name}}!"

	tests := []struct {
		name        string
		templateStr string
	}{
		{"Valid Template", validTemplateString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := New(tt.templateStr)
			if pt.template == nil {
				t.Errorf("expected template to be initialized, got nil")
			}
		})
	}
}

func TestNewPromptTemplateInvalidPanics(t *testing.T) {
	invalidTemplateString := "Hello, {{.name}"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for invalid template, but did not panic")
		}
	}()

	New(invalidTemplateString)
}

func TestFormat(t *testing.T) {
	tmplStr := "Hello, {{.Name}}!"
	pt := New(tmplStr)

	tests := []struct {
		name       string
		data       any
		want       string
		shouldFail bool
	}{
		{"Valid Data", map[string]string{"Name": "World"}, "Hello, World!", false},
		{"Missing Field", map[string]string{"WrongField": "World"}, "Hello, <no value>!", false},
		{"Nil Data", nil, "Hello, <no value>!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pt.Execute(tt.data)
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
