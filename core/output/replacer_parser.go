package output

import "strings"

type SymbolSeperatorParser struct{}

func (o *SymbolSeperatorParser) Parse(output string) (string, error) {
	var result strings.Builder
	for _, r := range output {
		if r == '?' || r == ',' || r == '.' || r == ';' || r == '!' {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}
	return result.String(), nil
}

func (p *SymbolSeperatorParser) ParseWithPrompt(output string, prompt PromptValue) (string, error) {
	// prompt not implemented.
	return output, nil
}

func (p *SymbolSeperatorParser) GetFormatInstructions() string {
	return "" // Eventually tell the large language model to return no symbols.
}
