package main

import (
	"fmt"

	"github.com/bit8bytes/gogantic/outputs/jsonout"
)

type joke struct {
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func main() {
	parser := jsonout.NewParser[joke]()
	joke, err := parser.Parse(`
	{ 
		"setup": "Why don't scientists trust atoms?", 
	 	"punchline": "Because they make up everything!"
	}`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Setup: " + joke.Setup)
	fmt.Println("Punchline: " + joke.Punchline)

	instructions := parser.Instructions()
	fmt.Println("LLM instructions: " + instructions)
}
