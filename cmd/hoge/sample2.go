package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func satisfyEvenInt2(x int) bool {
	return x%2 == 0
}

func analyzeValueSpecs(vv []*ast.ValueSpec, udt map[string]bool) {
	for _, v := range vv {
		if v.Values == nil || len(v.Values) == 0 {
			continue
		}
		ident, ok := v.Type.(*ast.Ident)
		if !ok {
			continue
		}
		if !udt[ident.Name] {
			continue
		}
		for _, value := range v.Values {
			lit, ok := value.(*ast.BasicLit)
			if !ok {
				continue
			}
			if lit.Kind != token.INT {
				continue
			}
			i, err := strconv.Atoi(lit.Value)
			if err != nil {
				continue
			}
			if !satisfyEvenInt2(i) {
				fmt.Println("value spec error: value is not even")
			}
		}
	}
}

func analyzeAssignStmts(as []*ast.AssignStmt, udt map[string]bool) {
	for _, a := range as {
		lhs, rhs := a.Lhs, a.Rhs
		// 一旦多重代入文は対応しない
		if len(lhs) != 1 || len(rhs) != 1 {
			fmt.Println("not support multiple assignment")
			continue
		}

		ident, ok := lhs[0].(*ast.Ident)
		if !ok {
			continue
		}
		if ident.Obj == nil {
			continue
		}
		decl, ok := ident.Obj.Decl.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if decl.Type == nil {
			continue
		}
		ident, ok = decl.Type.(*ast.Ident)
		if !ok {
			continue
		}
		if !udt[ident.Name] {
			continue
		}
		lit, ok := rhs[0].(*ast.BasicLit)
		if !ok {
			continue
		}
		i, err := strconv.Atoi(lit.Value)
		if err != nil {
			continue
		}
		if !satisfyEvenInt2(i) {
			fmt.Println("assign stmt error: value is not even")
		}
	}
}

func main() {

	typeFile := `package my_type
type Even int
`
	mainFile := `package main
func main() {
    var x Even = 10
	x = 10
    print(x)
}`
	fset := token.NewFileSet()
	tf, err := parser.ParseFile(fset, "type.go", typeFile, 0)
	if err != nil {
		panic("parse error" + err.Error())
	}

	// ユーザー定義型の名前を記録
	userDefineTypes := make(map[string]bool, 10) // TODO: 10は適当な値

	// 型定義をさがして見つける
	ast.Inspect(tf, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			userDefineTypes[typeSpec.Name.Name] = true // TODO: ここで型名を記録しているが、本来は名前だけでな型の制約も記述したい
		}
		return true
	})

	mf, err := parser.ParseFile(fset, "main.go", mainFile, 0)
	if err != nil {
		panic("parse error" + err.Error())
	}

	// 検査対象slice
	var targetValueSpecs []*ast.ValueSpec
	var targetAssignStmts []*ast.AssignStmt

	// main関数のファイルでast.ValueSpec, ast.AssignStmtを探して検査対象に追加する
	ast.Inspect(mf, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.ValueSpec:
			targetValueSpecs = append(targetValueSpecs, n.(*ast.ValueSpec))
		case *ast.AssignStmt:
			targetAssignStmts = append(targetAssignStmts, n.(*ast.AssignStmt))
		}
		return true
	})

	// 検査対象を一気に検査する
	analyzeValueSpecs(targetValueSpecs, userDefineTypes)
	analyzeAssignStmts(targetAssignStmts, userDefineTypes)

}
