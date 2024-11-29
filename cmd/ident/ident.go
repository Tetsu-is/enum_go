package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func satisfyOddInt(x int) bool {
	return x%2 == 1
}

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

	// まずはトップレベルの宣言からOdd型を取得する(ここでは取得成功をゴールとする)
	// Odd型はfile->Decls->GenDecl->Specs->TypeSpec->Nameに格納されている
	topLevelDecls := file.Decls
	for _, decl := range topLevelDecls {
		fmt.Println("top level decl")
		ast.Print(fset, decl)
		// 型宣言 TypeSpec だけを取り出したい
		// 型宣言はGenDeclで関数宣言はFuncDeclなので、GenDeclのみを取り出す
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

	}
	//ast.Print(fset, file)
	//ast.Print(fset, file.Scope.Objects["Odd"])
}
