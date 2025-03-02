package main

import (
	"fmt"

	"github.com/bit8bytes/gogantic/core/input"
)

func main() {
	fmt.Println("Example 1:")

	companyNamePrompt, _ := input.NewPromptTemplate("What is a good name for a company that makes {{.product}}?")
	data := map[string]string{"product": "coloful socks"}
	companyNameFormattedPrompt, _ := companyNamePrompt.Format(data)
	fmt.Println(companyNameFormattedPrompt)

	fmt.Println("Example 2:")

	buildProductPrompt, _ := input.NewPromptTemplate("{{.name}} want's to build {{.company}}.")
	buildProduct := map[string]string{"name": "Alex", "company": "coloful socks"}
	twoVariablesFormattedPrompt, _ := buildProductPrompt.Format(buildProduct)
	fmt.Println(twoVariablesFormattedPrompt)
}
