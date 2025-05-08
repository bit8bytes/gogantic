package main

import (
	"fmt"
	"os"

	"github.com/bit8bytes/gogantic/output"
	"github.com/bit8bytes/gogantic/output/markdown"
)

func main() {
	content, _ := os.ReadFile("markdown.md")

	var parser output.Parser[map[string]string] = &markdown.Parser{}
	parsedMarkdownOutput, _ := parser.Parse(string(content))

	section := "## Heading 1.1"
	if content, exists := parsedMarkdownOutput[section]; exists {
		fmt.Printf("Content of %s:\n%s...\n", section, content[:20]) // This will panic when using # Heading 1 because it is empty
	} else {
		fmt.Printf("Section %s not found.\n", section)
	}
}
