package replacer

import (
	"strings"

	"github.com/bit8bytes/gogantic/input/prompt"
)

type Symbols struct{}

func (o *Symbols) Parse(output string) (string, error) {
	var result strings.Builder
	for _, r := range output {
		if r == '?' || r == ',' || r == '.' || r == ';' || r == '!' {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}
	return result.String(), nil
}

func (p *Symbols) ParseWithPrompt(output string, prompt prompt.Prompt) (string, error) {
	// prompt not implemented.
	return output, nil
}

func (p *Symbols) GetFormatInstructions() string {
	return "" // Eventually tell the large language model to return no symbols.
}
