package output

import "github.com/bit8bytes/gogantic/input/prompt"

// Parser is an interface for parsing the output of an LLM call.
type Parser[T any] interface {
	Parse(text string) (T, error)
	ParseWithPrompt(text string, prompt prompt.Prompt) (T, error)
	GetFormatInstructions() string
}
