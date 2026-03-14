package tools

import (
	"context"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// GetFunctionSignature returns the signature and doc comment for a named function.
type GetFunctionSignature struct{}

func (t GetFunctionSignature) Name() string { return "GetFunctionSignature" }

func (t GetFunctionSignature) Description() string {
	return `Return the signature and doc comment for a named function or method.
Input: two lines — absolute path to a .go file or a directory containing .go files, then function name.
Use the working directory from the system prompt as the starting path.
Output: JSON object {name, receiver, params, results, doc}. First match wins.`
}

func (t GetFunctionSignature) Execute(ctx context.Context, input Input) (Output, error) {
	path, name, ok := twoLines(input.Content)
	if !ok {
		return Output{Content: `{"error":"provide path and function name on separate lines"}`}, nil
	}

	fset := token.NewFileSet()
	var files []*ast.File
	if strings.HasSuffix(path, ".go") {
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return Output{Content: `{"error":"` + err.Error() + `"}`}, nil
		}
		files = []*ast.File{f}
	} else {
		var err error
		files, err = walkGoFiles(path, fset)
		if err != nil {
			return Output{Content: `{"error":"` + err.Error() + `"}`}, nil
		}
	}

	type sig struct {
		Name     string `json:"name"`
		Receiver string `json:"receiver"`
		Params   string `json:"params"`
		Results  string `json:"results"`
		Doc      string `json:"doc"`
	}

	for _, f := range files {
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Name.Name != name {
				continue
			}
			recv := ""
			if fd.Recv != nil && len(fd.Recv.List) > 0 {
				recv = recvTypeName(fd.Recv.List[0].Type)
			}
			params, results := paramsAndResults(fd.Type)
			out, _ := json.Marshal(sig{fd.Name.Name, recv, params, results, docString(fd.Doc)})
			return Output{Content: string(out)}, nil
		}
	}

	return Output{Content: `{"error":"function not found"}`}, nil
}

func paramsAndResults(ft *ast.FuncType) (params, results string) {
	if ft == nil {
		return "()", ""
	}
	var ps []string
	if ft.Params != nil {
		for _, p := range ft.Params.List {
			ps = append(ps, exprString(p.Type))
		}
	}
	params = "(" + strings.Join(ps, ", ") + ")"

	if ft.Results != nil && len(ft.Results.List) > 0 {
		var rs []string
		for _, r := range ft.Results.List {
			rs = append(rs, exprString(r.Type))
		}
		results = "(" + strings.Join(rs, ", ") + ")"
	}
	return params, results
}
