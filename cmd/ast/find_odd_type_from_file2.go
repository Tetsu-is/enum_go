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

func registerTypeSpecs(file *ast.File, lVarList map[string]ast.TypeSpec) error {
	dcls := file.Decls
	for _, dcl := range dcls {
		genDcl, ok := dcl.(*ast.GenDecl)
		if !ok {
			continue
		}
		specs := genDcl.Specs
		for _, spec := range specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			typeName := typeSpec.Name.Name
			lVarList[typeName] = *typeSpec
		}
	}
	if len(lVarList) == 0 {
		return fmt.Errorf("Failed to find type declaration")
	}
	return nil
}

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

func getAssignStmts(mainFunc *ast.FuncDecl) ([]*ast.AssignStmt, error) {
	bodyList := mainFunc.Body.List
	var assignStmts []*ast.AssignStmt
	for _, stmt := range bodyList {
		asmtStmt, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}
		assignStmts = append(assignStmts, asmtStmt)
	}
	if len(assignStmts) == 0 {
		return nil, fmt.Errorf("Failed to find assignment statement")
	}
	return assignStmts, nil
}

func satisfyAssignStmt(assignStmt *ast.AssignStmt) error {
	lhs := assignStmt.Lhs
	lObj := lhs[0].(*ast.Ident).Obj
	lType := lObj.Decl.(*ast.ValueSpec).Type

	rhs := assignStmt.Rhs
	rhsBasicLit := rhs[0].(*ast.BasicLit)
	castedINTValue, err := strconv.Atoi(rhsBasicLit.Value)
	if err != nil {
		panic("Failed to cast value: " + err.Error())
	}

	isOdd := lType.(*ast.SelectorExpr).X.(*ast.Ident).Name == "my_type" && lType.(*ast.SelectorExpr).Sel.Name == "Odd"
	// if not Odd type, SKIP
	if !isOdd {
		return nil
	}

	// if not integer value, value does not satisfy Odd type
	if rhsBasicLit.Kind != token.INT {
		return fmt.Errorf("value is not integer type")
	}

	if satisfyOddInt(castedINTValue) {
		// msg tells place and reason of error
		return fmt.Errorf("satisfy error: %v/%d is not asssignable into Odd type", assignStmt.TokPos, castedINTValue)
	}

	return nil
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
	if err := registerTypeSpecs(tfile, lVarList); err != nil {
		panic("Failed to reg                                                                                                                ister type specs: " + err.Error())
	}
	// 2. parse main file & check if the variable is Odd type

	fPath := "./main.go"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fPath, nil, 0)
	if err != nil {
		panic("Failed to parse file: " + err.Error())
	}
	fmt.Println("Parsed main file")

	mainFunc, err := getMainFuncDecl(file)
	if err != nil {
		panic("Failed to find main function")
	}

	if mainFunc == nil {
		panic("Failed to find main function")
	}

	asmtStmts, err := getAssignStmts(mainFunc)
	var is_odd_type bool = false
	var is_odd_type_assigned bool = false
	// find assignment statement
	for _, asStmt := range asmtStmts {
		//init
		is_odd_type_assigned = false
		is_odd_type_assigned = false

		//ast.Print(fset, asmtStmt)

		lhs := asStmt.Lhs
		//fmt.Println("lhs")
		//ast.Print(fset, lhs)
		//fmt.Println("ast detail")
		//ast.Print(fset, lhs[0].(*ast.Ident).Obj.Decl.(*ast.ValueSpec).Names)
		lhsIdentObj := lhs[0].(*ast.Ident).Obj
		lhsType := lhsIdentObj.Decl.(*ast.ValueSpec).Type
		if lhsType.(*ast.SelectorExpr).X.(*ast.Ident).Name == "my_type" && lhsType.(*ast.SelectorExpr).Sel.Name == "Odd" {
			is_odd_type = true
		}
		rhs := asStmt.Rhs
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
			// error position ex. "maing.go:10"
			fmt.Println(fset.Position(asStmt.TokPos)) // fset.Positionを使えばfilename+lineができる
			fmt.Printf("%d is not assignable into Odd\n", castedIntValue)
		}
	}
	//assignStmts, err := getAssignStmts(mainFunc)
	//if err != nil {
	//	panic("failed to find assignment statement")
	//}

	//ast.Print(fset, mainBodyList)
	//ast.Print(fset, file)
}
