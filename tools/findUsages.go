package tools

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// FindUsages finds all occurrences of a named identifier in a Go package directory.
type FindUsages struct{}

func (t FindUsages) Name() string { return "FindUsages" }

func (t FindUsages) Description() string {
	return `Find all occurrences of a named identifier in a Go package directory.
AST-level only: includes declarations and usages — does not distinguish between them.
Input: two lines — absolute filesystem path to a directory containing .go files, then the symbol name.
Use the working directory from the system prompt as the starting path.
Output: one occurrence per line as file:line.`
}

func (t FindUsages) Execute(ctx context.Context, input Input) (Output, error) {
	dir, name, ok := twoLines(input.Content)
	if !ok {
		return Output{Content: `{"error":"provide directory and symbol name on separate lines"}`}, nil
	}

	fset := token.NewFileSet()
	files, err := walkGoFiles(dir, fset)
	if err != nil {
		return Output{Content: `{"error":"` + err.Error() + `"}`}, nil
	}

	var b strings.Builder
	found := 0
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if ok && ident.Name == name {
				pos := fset.Position(ident.Pos())
				fmt.Fprintf(&b, "%s:%d\n", pos.Filename, pos.Line)
				found++
			}
			return true
		})
	}

	if found == 0 {
		return Output{Content: "not found"}, nil
	}
	return Output{Content: b.String()}, nil
}
