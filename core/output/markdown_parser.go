package output

import "strings"

// Implemented Markdown Parsers
type MarkdownOutputParser struct{}

func (p *MarkdownOutputParser) Parse(output string) (map[string]string, error) {
	result := make(map[string]string)
	lines := strings.Split(output, "\n")

	var currentHeading string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			currentHeading = line
			result[currentHeading] = ""
		} else if currentHeading != "" {
			result[currentHeading] += line + "\n"
		}
	}

	for key, value := range result {
		result[key] = strings.TrimSpace(value)
	}
	return result, nil
}

func (p *MarkdownOutputParser) ParseWithPrompt(output string, prompt PromptValue) (map[string]string, error) {
	// prompt not implemented.
	return nil, nil
}

func (p *MarkdownOutputParser) GetFormatInstructions() string {
	return `
		Return the output as markdown in this format:

		# Title 1
		## Subtitle 1.1
		## Subtitle 1.2
		`
}
