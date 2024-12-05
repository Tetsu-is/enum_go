package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	src := `
	package main
	type Odd int
	func main() {
		var x Odd
		x = 1
		fmt.Println(x)
	}
	`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", src, 0)
	if err != nil {
		panic("Failed to parse file: " + err.Error())
	}
	ast.Print(fset, file)
	mainFunc := file.Decls[1].(*ast.FuncDecl)
	assignStmt := mainFunc.Body.List[1].(*ast.AssignStmt)
	fmt.Println("assignStmt")
	ast.Print(fset, assignStmt)
	ast.Print(fset, file.Decls[1].(*ast.FuncDecl).Body.List[1].(*ast.AssignStmt).Lhs[0].(*ast.Ident))
}
