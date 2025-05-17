package main

import (
	"fmt"

	"github.com/bit8bytes/gogantic/output"
	"github.com/bit8bytes/gogantic/output/json"
)

type Joke struct {
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func main() {
	var parser output.Parser[Joke] = &json.Parser[Joke]{}
	joke, err := parser.Parse(`{"setup": "Why don't scientists trust atoms?", "punchline": "Because they make up everything!"}`)
	if err != nil {
		panic(err)
	}

	fmt.Println(joke.Setup)
	fmt.Println(joke.Punchline)

	instructions := parser.GetFormatInstructions()
	fmt.Println(instructions)
}
