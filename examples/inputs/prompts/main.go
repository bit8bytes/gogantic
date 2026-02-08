package main

import (
	"fmt"

	"github.com/bit8bytes/gogantic/inputs/prompts"
)

func main() {
	type productData struct {
		Product string
	}

	companyNamePrompt := "What is a good name for a company that makes {{.Product}}?"
	data := productData{Product: "coloful socks"}
	companyNameFormattedPrompt, _ := prompts.New(companyNamePrompt).Execute(data)
	fmt.Println("Example 1: " + companyNameFormattedPrompt)

	type buildProductData struct {
		Name    string
		Company string
	}

	buildProductPrompt := "{{.Name}} want's to build {{.Company}}."
	buildProduct := buildProductData{Name: "Alex", Company: "coloful socks"}
	twoVariablesFormattedPrompt, _ := prompts.New(buildProductPrompt).Execute(buildProduct)
	fmt.Println("Example 2: " + twoVariablesFormattedPrompt)
}
