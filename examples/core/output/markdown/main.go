package main

import (
	"fmt"
	"os"

	"github.com/bit8bytes/gogantic/core/output"
)

func main() {
	content, _ := os.ReadFile("markdown.md")
	markdown := string(content)

	var parser output.OutputParser[map[string]string] = &output.MarkdownOutputParser{}
	parsedMarkdownOutput, _ := parser.Parse(markdown)

	section := "## Heading 1.1"
	if content, exists := parsedMarkdownOutput[section]; exists {
		fmt.Printf("Content of %s:\n%s...\n", section, content[:20]) // This will panic when using # Heading 1 because it is empty
	} else {
		fmt.Printf("Section %s not found.\n", section)
	}
}
