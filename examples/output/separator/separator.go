package main

import (
	"fmt"

	"github.com/bit8bytes/gogantic/output"
	"github.com/bit8bytes/gogantic/output/maps"
	"github.com/bit8bytes/gogantic/output/separator"
)

func main() {
	var parser output.Parser[string] = &separator.Text{}
	parsedStringOutput, _ := parser.Parse("hi, bye")
	fmt.Println(parsedStringOutput)

	var spaceParser output.Parser[[]string] = &separator.Space{}
	spaceResult, _ := spaceParser.Parse("example output with spaces")
	fmt.Println(spaceResult)

	var mapParser output.Parser[map[string]bool] = &maps.Parser{}
	mapResult, _ := mapParser.Parse("example output with spaces")
	fmt.Println(mapResult)

	var commaParser output.Parser[[]string] = &separator.Comma{}
	commaResult, _ := commaParser.Parse("example,output,with,commas")
	fmt.Println(commaResult)

	var sentenceParser output.Parser[[]string] = &separator.Sentence{}
	sentences, _ := sentenceParser.Parse("Hello! How are you? I am good.")
	fmt.Println(sentences)
}
