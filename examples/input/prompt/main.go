package main

import (
	"fmt"

	"github.com/bit8bytes/gogantic/input/prompt"
)

func main() {
	fmt.Println("Example 1:")

	companyNamePrompt := prompt.New("What is a good name for a company that makes {{.product}}?")
	data := map[string]string{"product": "coloful socks"}
	companyNameFormattedPrompt, _ := companyNamePrompt.Format(data)
	fmt.Println(companyNameFormattedPrompt)

	fmt.Println("Example 2:")

	buildProductPrompt := prompt.New("{{.name}} want's to build {{.company}}.")
	buildProduct := map[string]string{"name": "Alex", "company": "coloful socks"}
	twoVariablesFormattedPrompt, _ := buildProductPrompt.Format(buildProduct)
	fmt.Println(twoVariablesFormattedPrompt)
}
