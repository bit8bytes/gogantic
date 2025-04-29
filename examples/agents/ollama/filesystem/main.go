package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/bit8bytes/gogantic/agents"
	"github.com/bit8bytes/gogantic/core/models/llms/ollama"
)

// FileParams defines the structure for SaveToFile parameters
type FileParams struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
}

type OpenFile struct{}
type WriteAndSaveToFile struct{}

func main() {
	mistral_latest := ollama.OllamaModel{
		Model:   "gemma3:4b",
		Options: ollama.ModelOptions{NumCtx: 4096},
		Stream:  false,
		Stop:    []string{"\nObservation", "Observation"},
	}

	llm := ollama.New(mistral_latest)
	tools := map[string]agents.Tool{
		"OpenFile":           OpenFile{},
		"WriteAndSaveToFile": WriteAndSaveToFile{},
	}

	agent1 := agents.New(llm, tools)
	agent1.Task(`
1. Open the file foobar.txt. 
2. Read the content and add the sentence: I can edit files 
3. Save it to altered_foobar.txt
`)
	executor1 := agents.NewExecutor(agent1, agents.WithShowMessages())

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		executor1.Run(context.TODO())
		finalAnswer1, _ := agent1.GetFinalAnswer()
		fmt.Println("Agent 1 final answer:", finalAnswer1)
	}()

	wg.Wait()
	fmt.Println("Task completed")
}

func (c OpenFile) Name() string { return "OpenFile" }

func (c OpenFile) Call(ctx context.Context, input string) (string, error) {
	if input == "" {
		return "", errors.New(`please provide the filename in the Action Input: "filename"`)
	}

	content, err := os.ReadFile(input)
	if err != nil {
		return "", errors.New(err.Error())
	}

	response := "The content of the file is " + string(content)

	return response, nil
}

func (c WriteAndSaveToFile) Name() string { return "WriteAndSaveToFile" }

func (c WriteAndSaveToFile) Call(ctx context.Context, input string) (string, error) {
	fmt.Println(input)

	cleanedInput := strings.ReplaceAll(input, `\"`, `"`)
	cleanedInput = strings.ReplaceAll(cleanedInput, `'`, ``)
	cleanedInput = strings.Trim(cleanedInput, `"`)

	var params FileParams
	if err := json.Unmarshal([]byte(cleanedInput), &params); err != nil {
		fmt.Println("JSON parsing error:", err)

		if err := json.Unmarshal([]byte(input), &params); err != nil {
			return "", errors.New(`please provide a valid JSON with "content" and "filename" fields for tool "WriteAndSaveToFile"`)
		}
	}

	if params.Filename == "" {
		return "", errors.New(`please provide the filename in the "filename" field`)
	}

	if params.Content == "" {
		return "", errors.New(`please provide content to write in the "content" field`)
	}

	file, err := os.Create(params.Filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(params.Content)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully wrote to %s", params.Filename), nil
}
