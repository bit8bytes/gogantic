package output

import "strings"

// Implemented Map Parsers
type MapOutputParser struct{}

func (p *MapOutputParser) Parse(output string) (map[string]bool, error) {
	words := strings.Fields(output)
	wordMap := make(map[string]bool)
	for _, word := range words {
		wordMap[word] = true
	}
	return wordMap, nil
}

func (p *MapOutputParser) ParseWithPrompt(output string, prompt PromptValue) (map[string]bool, error) {
	// prompt not implemented.
	return nil, nil
}

func (p *MapOutputParser) GetFormatInstructions() string {
	return ""
}
