package jsonout

import (
	"testing"
)

func TestJsonout(t *testing.T) {
	type jokeData struct {
		Setup     string `json:"setup"`
		Punchline string `json:"punchline"`
	}

	parser := NewParser[jokeData]()
	joke, err := parser.Parse(`
	{ 
		"setup": "Why don't scientists trust atoms?", 
	 	"punchline": "Because they make up everything!"
	}`)
	if err != nil {
		panic(err)
	}

	if joke.Setup != "Why don't scientists trust atoms?" {
		t.Errorf("expected joke setup to be: %s", "Why don't scientists trust atoms?")
	}

	if joke.Punchline != "Because they make up everything!" {
		t.Errorf("expetec joke punchline to be: %s", "Because they make up everything!")
	}
}

func TestJsonoutArray(t *testing.T) {
	type jokeData struct {
		Setup     string `json:"setup"`
		Punchline string `json:"punchline"`
	}

	parser := NewParser[[]jokeData]()
	jokes, err := parser.Parse(`
	[
		{
			"setup": "Why don't scientists trust atoms?",
			"punchline": "Because they make up everything!"
		},
		{
			"setup": "What do you call a fake noodle?",
			"punchline": "An impasta!"
		}
	]`)
	if err != nil {
		t.Fatal(err)
	}

	if len(jokes) != 2 {
		t.Fatalf("expected 2 jokes, got %d", len(jokes))
	}

	if jokes[0].Setup != "Why don't scientists trust atoms?" {
		t.Errorf("expected first joke setup to be: %s, got: %s",
			"Why don't scientists trust atoms?", jokes[0].Setup)
	}

	if jokes[0].Punchline != "Because they make up everything!" {
		t.Errorf("expected first joke punchline to be: %s, got: %s",
			"Because they make up everything!", jokes[0].Punchline)
	}

	if jokes[1].Setup != "What do you call a fake noodle?" {
		t.Errorf("expected second joke setup to be: %s, got: %s",
			"What do you call a fake noodle?", jokes[1].Setup)
	}

	if jokes[1].Punchline != "An impasta!" {
		t.Errorf("expected second joke punchline to be: %s, got: %s",
			"An impasta!", jokes[1].Punchline)
	}
}
