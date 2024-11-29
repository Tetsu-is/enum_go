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
	ast.Print(fset, file)
	// 続いてAssignStmtでOdd型に偶数値を代入していないかをチェックする
	// AssignStmtはast.Function->Body->List->ast.Assignmentに格納されている
	// まずはmain関数を取得する
	for _, genDecl := range file.Decls {
		funcDecl, ok := genDecl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Name.Name != "main" {
			continue
		}

		// main関数の中を探す
		for _, stmt := range funcDecl.Body.List {
			assignStmt, ok := stmt.(*ast.AssignStmt)
			if !ok {
				continue
			}

			// AssignStmtの左辺にOdd型があるかをチェックする
			for _, lhs := range assignStmt.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok {
					continue
				}
				if ident.Name == "x" { // Odd型として定義されている変数名x
					fmt.Println("Found x")
					ast.Print(fset, assignStmt)
				}
			}
		}
	}
	ast.Print(fset, file)

}
