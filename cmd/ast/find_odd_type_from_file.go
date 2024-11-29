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
	fPath := "./my_type/odd.go"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fPath, nil, 0)
	if err != nil {
		panic("Failed to parse file: " + err.Error())
	}

	// まずはトップレベルの宣言からOdd型を取得する(ここでは取得成功をゴールとする)
	// Odd型はfile->Decls->GenDecl->Specs->TypeSpec->Nameに格納されている
	topLevelDecls := file.Decls
	for _, decl := range topLevelDecls {
		// 型宣言 TypeSpec だけを取り出したい
		// 型宣言はGenDeclで関数宣言はFuncDeclなので、GenDeclのみを取り出す
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		specs := genDecl.Specs // この中にTypeSpecがある
		for _, spec := range specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// if not Odd type, SKIP
			if typeSpec.Name.Name != "Odd" {
				continue
			}

			fmt.Println("Found Odd type")
			ast.Print(fset, typeSpec)
		}

	}
	//ast.Print(fset, file)
	//ast.Print(fset, file.Scope.Objects["Odd"])
}
