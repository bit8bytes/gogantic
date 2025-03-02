package main

import (
	"fmt"

	"github.com/bit8bytes/gogantic/core/output"
)

func main() {
	var parser output.OutputParser[string] = &output.StringOutputParser{}
	parsedStringOutput, _ := parser.Parse("hi, bye")
	fmt.Println(parsedStringOutput)

	var spaceParser output.OutputParser[[]string] = &output.SpaceSeparatedListOutputParser{}
	spaceResult, _ := spaceParser.Parse("example output with spaces")
	fmt.Println(spaceResult)

	var mapParser output.OutputParser[map[string]bool] = &output.MapOutputParser{}
	mapResult, _ := mapParser.Parse("example output with spaces")
	fmt.Println(mapResult)

	var commaParser output.OutputParser[[]string] = &output.CommaSeparatedListOutputParser{}
	commaResult, _ := commaParser.Parse("example,output,with,commas")
	fmt.Println(commaResult)

	var sentenceParser output.OutputParser[[]string] = &output.SentenceSeparatedListOutputParser{}
	sentences, _ := sentenceParser.Parse("Hello! How are you? I am good.")
	fmt.Println(sentences)
}
