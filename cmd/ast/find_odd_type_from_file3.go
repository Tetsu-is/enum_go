package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func getMainFuncDecl(file *ast.File) (*ast.FuncDecl, error) {
	dcls := file.Decls
	var mainFuncDecl *ast.FuncDecl
	for _, dcl := range dcls {
		fdcl, ok := dcl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fdcl.Name.Name != "main" {
			continue
		}
		mainFuncDecl = fdcl
	}
	if mainFuncDecl == nil {
		return nil, fmt.Errorf("failed to find main function")
	}
	return mainFuncDecl, nil
}

func main() {
	// TODO: var declの個数で格納領域を決定できたら良さそう
	lVarList := make(map[string]ast.TypeSpec, 0)

	// parse type file
	typeFilePath := "./main.go"
	typeFset := token.NewFileSet()
	typeFile, err := parser.ParseFile(typeFset, typeFilePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	mainFilePath := "./main.go"
	mainFset := token.NewFileSet()
	mainFile, err := parser.ParseFile(mainFset, mainFilePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	mainFunc, err := getMainFuncDecl(mainFile)
	if err != nil {
		panic(err)
	}

}
