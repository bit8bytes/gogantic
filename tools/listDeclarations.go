package tools

import (
	"context"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// ListDeclarations lists all top-level declarations in a Go file or package directory.
type ListDeclarations struct{}

func (t ListDeclarations) Name() string { return "ListDeclarations" }

func (t ListDeclarations) Description() string {
	return `List all top-level declarations (funcs, types, vars, consts) in a Go file or package directory.
Input: absolute path to a .go file or a directory containing .go files.
Use the working directory from the system prompt as the starting path.
Output: JSON array of {kind, name, doc}.`
}

func (t ListDeclarations) Execute(ctx context.Context, input Input) (Output, error) {
	path := strings.TrimSpace(input.Content)
	if path == "" {
		return Output{Content: `{"error":"provide a file or directory path"}`}, nil
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

	type decl struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
		Doc  string `json:"doc"`
	}

	var decls []decl
	for _, f := range files {
		for _, d := range f.Decls {
			switch fd := d.(type) {
			case *ast.FuncDecl:
				decls = append(decls, decl{"func", fd.Name.Name, docString(fd.Doc)})
			case *ast.GenDecl:
				for _, spec := range fd.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						decls = append(decls, decl{"type", s.Name.Name, docString(fd.Doc)})
					case *ast.ValueSpec:
						kind := fd.Tok.String()
						for _, name := range s.Names {
							decls = append(decls, decl{kind, name.Name, docString(fd.Doc)})
						}
					}
				}
			}
		}
	}

	out, _ := json.Marshal(decls)
	return Output{Content: string(out)}, nil
}
