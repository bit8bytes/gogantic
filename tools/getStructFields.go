package tools

import (
	"context"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// GetStructFields returns the fields, types, and struct tags for a named struct type.
type GetStructFields struct{}

func (t GetStructFields) Name() string { return "GetStructFields" }

func (t GetStructFields) Description() string {
	return `Return all fields, their types, and struct tags for a named struct type.
Input: two lines — absolute path to a .go file or a directory containing .go files, then struct name.
Use the working directory from the system prompt as the starting path.
Output: JSON array of {name, type, tag}.`
}

func (t GetStructFields) Execute(ctx context.Context, input Input) (Output, error) {
	path, name, ok := twoLines(input.Content)
	if !ok {
		return Output{Content: `{"error":"provide path and struct name on separate lines"}`}, nil
	}

	fset := token.NewFileSet()
	var files []*ast.File
	if strings.HasSuffix(path, ".go") {
		f, err := parser.ParseFile(fset, path, nil, 0)
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

	type field struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Tag  string `json:"tag"`
	}

	for _, f := range files {
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.Name != name {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return Output{Content: `{"error":"` + name + ` is not a struct"}`}, nil
				}
				var fields []field
				for _, f := range st.Fields.List {
					typ := exprString(f.Type)
					tag := ""
					if f.Tag != nil {
						tag = strings.Trim(f.Tag.Value, "`")
					}
					if len(f.Names) == 0 {
						// embedded field
						fields = append(fields, field{typ, typ, tag})
					}
					for _, n := range f.Names {
						fields = append(fields, field{n.Name, typ, tag})
					}
				}
				out, _ := json.Marshal(fields)
				return Output{Content: string(out)}, nil
			}
		}
	}

	return Output{Content: `{"error":"struct not found"}`}, nil
}
