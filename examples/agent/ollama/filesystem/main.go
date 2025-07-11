package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bit8bytes/gogantic/agent"
	"github.com/bit8bytes/gogantic/llm/ollama"
	"github.com/bit8bytes/gogantic/runner"
	"github.com/bit8bytes/gogantic/tool"
)

// FileParams defines the structure for SaveToFile parameters
type FileParams struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
}

type OpenFile struct{}
type WriteAndSaveToFile struct{}

func main() {
	model := ollama.Model{
		Model:   "gemma3:4b",
		Options: ollama.Options{NumCtx: 4096},
		Stream:  false,
		Stop:    []string{"\nObservation", "Observation"},
	}

	llm := ollama.New(model)
	tools := map[string]tool.Tool{
		"OpenFile":           OpenFile{},
		"WriteAndSaveToFile": WriteAndSaveToFile{},
	}

	agent := agent.New(llm, tools)
	agent.Task(`
1. Open the file foobar.txt. 
2. Read the content and add the sentence: I can edit files. I am a happy local Agent!
3. Save it to altered_foobar.txt
`)
	runner := runner.New(agent, runner.WithShowMessages())

	runner.Run(context.TODO())
	finalAnswer1, _ := agent.GetFinalAnswer()
	fmt.Println("Agent 1 final answer:", finalAnswer1)
}

func (c OpenFile) Name() string { return "OpenFile" }

func (c OpenFile) Call(ctx context.Context, input tool.Input) (tool.Output, error) {
	if input.Content == "" {
		return tool.Output{Content: ""}, errors.New(`please provide the filename in the Action Input: "filename"`)
	}

	content, err := os.ReadFile(input.Content)
	if err != nil {
		return tool.Output{Content: ""}, errors.New(err.Error())
	}

	response := "The content of the file is " + string(content)

	return tool.Output{Content: response}, nil
}

func (c WriteAndSaveToFile) Name() string { return "WriteAndSaveToFile" }

func (c WriteAndSaveToFile) Call(ctx context.Context, input tool.Input) (tool.Output, error) {
	fmt.Println(input)

	cleanedInput := strings.ReplaceAll(input.Content, `\"`, `"`)
	cleanedInput = strings.ReplaceAll(cleanedInput, `'`, ``)
	cleanedInput = strings.Trim(cleanedInput, `"`)

	var params FileParams
	if err := json.Unmarshal([]byte(cleanedInput), &params); err != nil {
		fmt.Println("JSON parsing error:", err)

		if err := json.Unmarshal([]byte(input.Content), &params); err != nil {
			return tool.Output{}, errors.New(`please provide a valid JSON with "content" and "filename" fields for tool "WriteAndSaveToFile"`)
		}
	}

	if params.Filename == "" {
		return tool.Output{}, errors.New(`please provide the filename in the "filename" field`)
	}

	if params.Content == "" {
		return tool.Output{}, errors.New(`please provide content to write in the "content" field`)
	}

	file, err := os.Create(params.Filename)
	if err != nil {
		return tool.Output{}, err
	}
	defer file.Close()

	_, err = file.WriteString(params.Content)
	if err != nil {
		return tool.Output{}, err
	}

	return tool.Output{Content: fmt.Sprintf("Successfully wrote to %s", params.Filename)}, nil
}
