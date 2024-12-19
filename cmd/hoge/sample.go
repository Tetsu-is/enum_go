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

type Enum struct {
	Package string
	Name    string
	Type    string
}

func registerEnum(vs ast.ValueSpec, m map[string]Enum) {
	// pick up the value
	n := vs.Names[0].Name
	// pick up type
	t := vs.Type.(*ast.Ident).Name

	m[n] = Enum{
		Package: "hoge",
		Name:    n,
		Type:    t,
	}
}

func main() {
	typeDef := `package my_type
type Even int
const Two Even = 2
const Four Even = 4
`
	src1 := `package main
func main() {
    var x Even = 10
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
	// TODO: type <-> satisfy関数の対応もつけたいよね
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

	ast.Print(fset, file)

	// リテラル代入を検索
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ValueSpec: // var宣言
			if node.Values != nil && len(node.Values) > 0 {
				ident, ok := node.Type.(*ast.Ident)
				if !ok {
					return false
				}
				if !userDefinedTypes[ident.Name] {
					return false
				}
				for _, value := range node.Values {
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
					if !satisfyEvenInt(i) {
						fmt.Printf("value is not even at %v\n", fset.Position(node.Pos()))
					}
				}
			}
		case *ast.AssignStmt: // 代入文
			for i, lhs := range node.Lhs {
				ident, ok := lhs.(*ast.Ident)
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
				if !userDefinedTypes[ident.Name] {
					continue
				}
				lit, ok := node.Rhs[i].(*ast.BasicLit)
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
				if !satisfyEvenInt(i) {
					fmt.Printf("value is not even at %v\n", fset.Position(node.Pos()))
				}
			}
		}
		return true
	})
}
