// Command check-go-docs verifies that every exported Go declaration has a
// documentation comment suitable for GoDoc.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: check-go-docs <package-directory>...")
		os.Exit(2)
	}

	var missing []string
	for _, directory := range os.Args[1:] {
		missing = append(missing, undocumented(directory)...)
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		fmt.Fprintln(os.Stderr, "Exported Go declarations need documentation:")
		for _, declaration := range missing {
			fmt.Fprintf(os.Stderr, "  %s\n", declaration)
		}
		os.Exit(1)
	}
}

func undocumented(directory string) []string {
	files, err := parser.ParseDir(token.NewFileSet(), directory, func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", directory, err)
		os.Exit(1)
	}

	var missing []string
	for _, packageFiles := range files {
		for filename, file := range packageFiles.Files {
			for _, declaration := range file.Decls {
				switch declaration := declaration.(type) {
				case *ast.FuncDecl:
					if exportedFunction(declaration) && declaration.Doc == nil {
						missing = append(missing, location(filename, declaration.Name.Name))
					}
				case *ast.GenDecl:
					for _, specification := range declaration.Specs {
						switch specification := specification.(type) {
						case *ast.TypeSpec:
							if specification.Name.IsExported() && declaration.Doc == nil && specification.Doc == nil {
								missing = append(missing, location(filename, specification.Name.Name))
							}
						case *ast.ValueSpec:
							for _, name := range specification.Names {
								if name.IsExported() && declaration.Doc == nil && specification.Doc == nil {
									missing = append(missing, location(filename, name.Name))
								}
							}
						}
					}
				}
			}
		}
	}
	return missing
}

func exportedFunction(declaration *ast.FuncDecl) bool {
	if !declaration.Name.IsExported() {
		return false
	}
	if declaration.Recv == nil {
		return true
	}
	if len(declaration.Recv.List) != 1 {
		return false
	}
	receiver := declaration.Recv.List[0].Type
	if pointer, ok := receiver.(*ast.StarExpr); ok {
		receiver = pointer.X
	}
	name, ok := receiver.(*ast.Ident)
	return ok && name.IsExported()
}

func location(filename, name string) string {
	return filepath.ToSlash(filename) + ": " + name
}
