package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

type Driver struct {
	Orders       int
	DrivingYears int
}

func main() {
	m := map[string]int64{"orders": 100000, "driving_years": 4}
	rule := `orders > 10000 && driving_years > 5`
	fmt.Println(Eval(m, rule))
}

func Eval(m map[string]int64, expr string) (bool, error) {
	exprAst, err := parser.ParseExpr(expr)
	if err != nil {
		return false, err
	}

	fset := token.NewFileSet()
	ast.Print(fset, exprAst)
	return judge(exprAst, m), nil
}

func judge(bop ast.Node, m map[string]int64) bool {
	// 叶子节点
	if isLeaf(bop) {
		expr := bop.(*ast.BinaryExpr)
		x := expr.X.(*ast.Ident)
		y := expr.Y.(*ast.BasicLit)

		if expr.Op == token.GTR {
			left := m[x.Name]
			right, _ := strconv.ParseInt(y.Value, 10, 64)
			return left > right
		}
		return false
	}

	expr, ok := bop.(*ast.BinaryExpr)
	if !ok {
		println("this cannot be true")
		return false
	}

	switch expr.Op {
	case token.LAND:
		return judge(expr.X, m) && judge(expr.Y, m)
	case token.LOR:
		return judge(expr.X, m) || judge(expr.Y, m)
	}

	println("unsupported operator")
	return false
}

func isLeaf(bop ast.Node) bool {
	expr, ok := bop.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	// 二元表达式的最小单位，左节点是标识符，右节点是值
	_, okL := expr.X.(*ast.Ident)
	_, okR := expr.Y.(*ast.BasicLit)
	if okL && okR {
		return true
	}

	return false
}
