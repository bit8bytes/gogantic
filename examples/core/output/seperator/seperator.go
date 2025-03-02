package main

import (
	"fmt"
	"strings"

	"github.com/bit8bytes/gogantic/core/output"
)

func main() {
	modelOutput := "Hello, world! How are you? I'm fine; thanks."

	// This is useful running before the whitelist/blacklist evaluator
	// The evaluator splits the string into words which could contain symbols and words such as "Hello,"
	var seperator output.OutputParser[string] = &output.SymbolSeperatorParser{}
	seperatedResult, _ := seperator.Parse(modelOutput)
	fmt.Println(seperatedResult)

	// Example blacklist: Semicolon (';') is not allowed
	if strings.Contains(seperatedResult, ";") {
		fmt.Println("Example blacklist: Semicolon (';') is not allowed")
	}
}
