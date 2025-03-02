package output

import (
	"regexp"
	"strings"
)

// Implemented String Parsers
type StringOutputParser struct{}
type SpaceSeparatedListOutputParser struct{}
type CommaSeparatedListOutputParser struct{}
type SentenceSeparatedListOutputParser struct{}

func (p *StringOutputParser) Parse(output string) (string, error) {
	return output, nil
}

func (p *StringOutputParser) ParseWithPrompt(output string, prompt PromptValue) (string, error) {
	// prompt not implemented.
	return output, nil
}

func (p *StringOutputParser) GetFormatInstructions() string {
	return "Return the output in the requested format."
}

func (p *SpaceSeparatedListOutputParser) Parse(output string) ([]string, error) {
	output = strings.TrimSpace(output)
	return strings.Split(output, " "), nil
}

func (p *SpaceSeparatedListOutputParser) ParseWithPrompt(output string, prompt PromptValue) ([]string, error) {
	// prompt not implemented.
	return strings.Split(output, " "), nil
}

func (p *SpaceSeparatedListOutputParser) GetFormatInstructions() string {
	return "Output only the result."
}

func (p *CommaSeparatedListOutputParser) Parse(output string) ([]string, error) {
	output = strings.TrimSpace(output)
	return strings.Split(output, ","), nil
}

func (p *CommaSeparatedListOutputParser) ParseWithPrompt(output string, prompt PromptValue) ([]string, error) {
	// prompt not implemented.
	return strings.Split(output, ","), nil
}

func (p *CommaSeparatedListOutputParser) GetFormatInstructions() string {
	return "Output only the result."
}

func (p *SentenceSeparatedListOutputParser) Parse(output string) ([]string, error) {
	re := regexp.MustCompile(`[?.!]\s*`)
	sentences := re.Split(output, -1)
	for i, sentence := range sentences {
		sentences[i] = strings.TrimSpace(sentence)
	}
	return sentences, nil
}

func (p *SentenceSeparatedListOutputParser) ParseWithPrompt(output string, prompt PromptValue) ([]string, error) {
	re := regexp.MustCompile(`[?.!]\s*`)
	sentences := re.Split(output, -1)
	for i, sentence := range sentences {
		sentences[i] = strings.TrimSpace(sentence)
	}
	return sentences, nil
}

func (p *SentenceSeparatedListOutputParser) GetFormatInstructions() string {
	return ""
}
