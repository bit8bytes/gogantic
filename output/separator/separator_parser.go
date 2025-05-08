package separator

import (
	"regexp"
	"strings"

	"github.com/bit8bytes/gogantic/input/prompt"
)

// Implemented String Parsers
type Text struct{}
type Space struct{}
type Comma struct{}
type Sentence struct{}

func (p *Text) Parse(output string) (string, error) {
	return output, nil
}

func (p *Text) ParseWithPrompt(output string, prompt prompt.Prompt) (string, error) {
	// prompt not implemented.
	return output, nil
}

func (p *Text) GetFormatInstructions() string {
	return "Return the output in the requested format."
}

func (p *Space) Parse(output string) ([]string, error) {
	output = strings.TrimSpace(output)
	return strings.Split(output, " "), nil
}

func (p *Space) ParseWithPrompt(output string, prompt prompt.Prompt) ([]string, error) {
	// prompt not implemented.
	return strings.Split(output, " "), nil
}

func (p *Space) GetFormatInstructions() string {
	return "Output only the result."
}

func (p *Comma) Parse(output string) ([]string, error) {
	output = strings.TrimSpace(output)
	return strings.Split(output, ","), nil
}

func (p *Comma) ParseWithPrompt(output string, prompt prompt.Prompt) ([]string, error) {
	// prompt not implemented.
	return strings.Split(output, ","), nil
}

func (p *Comma) GetFormatInstructions() string {
	return "Output only the result."
}

func (p *Sentence) Parse(output string) ([]string, error) {
	re := regexp.MustCompile(`[?.!]\s*`)
	sentences := re.Split(output, -1)
	for i, sentence := range sentences {
		sentences[i] = strings.TrimSpace(sentence)
	}
	return sentences, nil
}

func (p *Sentence) ParseWithPrompt(output string, prompt prompt.Prompt) ([]string, error) {
	re := regexp.MustCompile(`[?.!]\s*`)
	sentences := re.Split(output, -1)
	for i, sentence := range sentences {
		sentences[i] = strings.TrimSpace(sentence)
	}
	return sentences, nil
}

func (p *Sentence) GetFormatInstructions() string {
	return ""
}
