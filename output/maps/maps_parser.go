package maps

import (
	"strings"

	"github.com/bit8bytes/gogantic/input/prompt"
)

// Implemented Map Parsers
type Parser struct{}

func (p *Parser) Parse(output string) (map[string]bool, error) {
	words := strings.Fields(output)
	wordMap := make(map[string]bool)
	for _, word := range words {
		wordMap[word] = true
	}
	return wordMap, nil
}

func (p *Parser) ParseWithPrompt(output string, prompt prompt.Prompt) (map[string]bool, error) {
	// prompt not implemented.
	return nil, nil
}

func (p *Parser) GetFormatInstructions() string {
	return ""
}
