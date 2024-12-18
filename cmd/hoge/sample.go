package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func satisfyEvenInt(x int) bool {
	return x%2 == 0
}

func main() {
	typeDef := `package my_type
type Even int`
	src1 := `package main
func main() {
    var x Even = 0
	x = 10
    print(x)
}`

	// 型定義の解析
	fset := token.NewFileSet()
	tf, err := parser.ParseFile(fset, "type.go", typeDef, 0)
	if err != nil {
		panic("parse error" + err.Error())
	}

	// ユーザー定義型の名前を記録
	userDefinedTypes := make(map[string]bool)
	ast.Inspect(tf, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			userDefinedTypes[typeSpec.Name.Name] = true
		}
		return true
	})

	// メインコードの解析
	file, err := parser.ParseFile(fset, "main.go", src1, 0)
	if err != nil {
		panic("parse error" + err.Error())
	}

	// リテラル代入を検索
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ValueSpec: // var宣言
			if node.Values != nil && len(node.Values) > 0 {
				if ident, ok := node.Type.(*ast.Ident); ok {
					if userDefinedTypes[ident.Name] {
						for _, value := range node.Values {
							if lit, ok := value.(*ast.BasicLit); ok {
								if lit.Kind == token.INT {
									if i, err := strconv.Atoi(lit.Value); err == nil {
										if !satisfyEvenInt(i) {
											fmt.Printf("value is not even at %v\n", fset.Position(node.Pos()))
										}
									}
								}
							}
						}
					}
				}
			}
		case *ast.AssignStmt: // 代入文
			for i, lhs := range node.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if ident.Obj != nil {
						if decl, ok := ident.Obj.Decl.(*ast.ValueSpec); ok {
							if decl.Type != nil {
								if ident, ok := decl.Type.(*ast.Ident); ok {
									if userDefinedTypes[ident.Name] {
										if lit, ok := node.Rhs[i].(*ast.BasicLit); ok {
											if lit.Kind == token.INT {
												if i, err := strconv.Atoi(lit.Value); err == nil {
													if !satisfyEvenInt(i) {
														fmt.Printf("value is not even at %v\n", fset.Position(node.Pos()))
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return true
	})
}
