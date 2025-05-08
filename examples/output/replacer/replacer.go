package main

import (
	"fmt"
	"strings"

	"github.com/bit8bytes/gogantic/output"
	"github.com/bit8bytes/gogantic/output/replacer"
)

func main() {
	modelOutput := "Hello, world! How are you? I'm fine; thanks."

	// This is useful running before the whitelist/blacklist evaluator
	// The evaluator splits the string into words which could contain symbols and words such as "Hello,"
	var seperator output.Parser[string] = &replacer.Symbols{}
	seperatedResult, _ := seperator.Parse(modelOutput)
	fmt.Println(seperatedResult)

	// Example blacklist: Semicolon (';') is not allowed
	if strings.Contains(seperatedResult, ";") {
		fmt.Println("Example blacklist: Semicolon (';') is not allowed")
	}
}
