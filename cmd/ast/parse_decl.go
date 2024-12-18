package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	src := `
		package main
		
		func main() {
			var x my_type.Odd =1
			print(x)
		}
	`

	demo := ast.ValueSpec{
		Names:   []*ast.Ident{
			{
				NamePos: token.Pos(0),
				Name:   "x",
				Obj: nil,
			},
		},
		Type:    *ast.SelectorExpr{
			X: *ast.Ident{
				NamePos: token.Pos(0),
				Name:    "my_type",
			},
			Sel: *ast.Ident{
				NamePos: token.Pos(0),
				Name:    "Odd",
			},
		}
		Values:  nil,
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", src, 0)
	if err != nil {
		panic("failed to parse file: " + err.Error())
	}

	ast.Print(fset, file)
}
