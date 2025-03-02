package output 

type PromptValue struct {
    Content string
}

// OutputParser is an interface for parsing the output of an LLM call.
type OutputParser[T any] interface {
	Parse(text string) (T, error)
	ParseWithPrompt(text string, prompt PromptValue) (T, error)
	GetFormatInstructions() string
}