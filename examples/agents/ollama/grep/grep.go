package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bit8bytes/gogantic/agents/tools"
)

type ListDir struct{}

func (t ListDir) Name() string { return "ListDir" }
func (t ListDir) Description() string {
	return "List all files and folders in a directory. Input: a directory path (e.g. 'examples/'). Returns one entry per line."
}

func (t ListDir) Execute(ctx context.Context, input tools.Input) (tools.Output, error) {
	dir := strings.TrimSpace(input.Content)
	if dir == "" {
		return tools.Output{Content: "usage: <directory>"}, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return tools.Output{}, err
	}

	var dirs, files []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name()+"/")
		} else {
			files = append(files, entry.Name())
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Directory: %s\n", dir)
	fmt.Fprintf(&b, "Folders: %s\n", strings.Join(dirs, ", "))
	fmt.Fprintf(&b, "Files: %s", strings.Join(files, ", "))
	return tools.Output{Content: b.String()}, nil
}
