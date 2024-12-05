package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func satisfyOddInt(x int) bool {
	return x%2 == 1
}

func main() {
	lVarList := map[string]ast.TypeSpec{}

	// 1. parse type definition file & register variables into lVarList
	tfPath := "./my_type/odd.go"
	tfset := token.NewFileSet()
	tfile, err := parser.ParseFile(tfset, tfPath, nil, 0)
	if err != nil {
		panic("Failed to parse file: " + err.Error())
	}
	fmt.Println("Parsed type definition file")
	ast.Print(tfset, tfile)
	// scan top level declarations
	topLevelDecls := tfile.Decls
	for _, decl := range topLevelDecls {
		// pick up only type declaration
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		specs := genDecl.Specs // this contains TypeSpec
		for _, spec := range specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			// if it is type declaration, register the variable into lVarList
			typeName := typeSpec.Name.Name
			lVarList[typeName] = *typeSpec
			fmt.Println("Save type declaration: " + typeName)
			ast.Print(tfset, lVarList[typeName])
		}
	}
	// 2. parse main file & check if the variable is Odd type

	fPath := "./main.go"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fPath, nil, 0)
	if err != nil {
		panic("Failed to parse file: " + err.Error())
	}
	fmt.Println("Parsed main file")

	// find assignment statement
	decls := file.Decls
	var mainFunc *ast.FuncDecl
	for _, d := range decls {
		// find main function
		funcDecl, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Name.Name != "main" {
			continue
		}

		mainFunc = funcDecl
	}

	if mainFunc == nil {
		panic("Failed to find main function")
	}

	mainBodyList := mainFunc.Body.List
	var is_odd_type bool = false
	var is_odd_type_assigned bool = false
	// find assignment statement
	for _, stmt := range mainBodyList {
		//init
		is_odd_type_assigned = false
		is_odd_type_assigned = false

		asmtStmt, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}
		fmt.Println("Found assignment statement")
		//ast.Print(fset, asmtStmt)

		lhs := asmtStmt.Lhs
		//fmt.Println("lhs")
		//ast.Print(fset, lhs)
		//fmt.Println("ast detail")
		//ast.Print(fset, lhs[0].(*ast.Ident).Obj.Decl.(*ast.ValueSpec).Names)
		lhsIdentObj := lhs[0].(*ast.Ident).Obj
		lhsType := lhsIdentObj.Decl.(*ast.ValueSpec).Type
		if lhsType.(*ast.SelectorExpr).X.(*ast.Ident).Name == "my_type" && lhsType.(*ast.SelectorExpr).Sel.Name == "Odd" {
			is_odd_type = true
		}
		rhs := asmtStmt.Rhs
		//fmt.Println("rhs")
		rhsBasicLit := rhs[0].(*ast.BasicLit)
		castedIntValue, err := strconv.Atoi(rhsBasicLit.Value)
		if err != nil {
			panic("Failed to cast value: " + err.Error())
		}

		if is_odd_type && rhsBasicLit.Kind == token.INT && satisfyOddInt(castedIntValue) {
			is_odd_type_assigned = true
		}
		//ast.Print(fset, rhs)
		if is_odd_type && !is_odd_type_assigned {
			fmt.Printf("%d is not assignable into Odd\n", castedIntValue)
		}
	}
	//ast.Print(fset, mainBodyList)
	//ast.Print(fset, file)
}
