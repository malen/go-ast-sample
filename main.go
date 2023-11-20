package main

// https://www.cnblogs.com/apperception/p/16399821.html
import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"

	"github.com/Knetic/govaluate"
)

func main() {
	m := []map[string]interface{}{
		{"name": "aoi", "age": 18, "money": 100.2},
		{"name": "aoo", "age": 20, "money": 10.2},
		{"name": "aon", "age": 18, "money": 58.6},
	}
	rule := `(age >= 20 || name == "aoi") && money > 100`
	for _, v := range m {
		result, _ := Eval(v, rule)
		println(result)
	}
}

func Eval(m map[string]interface{}, expr string) (bool, error) {
	exprAst, err := parser.ParseExpr(expr)
	if err != nil {
		return false, err
	}

	fset := token.NewFileSet()
	ast.Print(fset, exprAst)
	return judge(exprAst, m), nil
}

func judge(bop ast.Node, m map[string]interface{}) bool {

	if isLeaf(bop) {
		expr := bop.(*ast.BinaryExpr)
		// 类型断言
		x := expr.X.(*ast.Ident)
		y := expr.Y.(*ast.BasicLit)

		var evalExpr *govaluate.EvaluableExpression
		switch t := m[x.Name].(type) {
		case string:
			evalExpr, _ = govaluate.NewEvaluableExpression(fmt.Sprintf(`"%s" %s %s`, t, expr.Op.String(), y.Value))
		case int:
			right, _ := strconv.ParseInt(y.Value, 10, 64)
			evalExpr, _ = govaluate.NewEvaluableExpression(fmt.Sprintf("%d %s %d", t, expr.Op.String(), right))
		case float64:
			right, _ := strconv.ParseFloat(y.Value, 64)
			evalExpr, _ = govaluate.NewEvaluableExpression(fmt.Sprintf("%f %s %f", t, expr.Op.String(), right))
		default:
		}

		result, _ := evalExpr.Evaluate(nil)

		// if expr.Op == token.GTR {
		// 	left := m[x.Name]
		// 	right, _ := strconv.ParseInt(y.Value, 10, 64)
		// 	return left > right
		// } else if expr.Op == token.GEQ {
		// 	left := m[x.Name]
		// 	right, _ := strconv.ParseInt(y.Value, 10, 64)
		// 	return left >= right
		// }
		m1 := fmt.Sprint(result)
		r, _ := strconv.ParseBool(m1)
		return r
	}

	switch guess := bop.(type) {
	case *ast.ParenExpr:
		println("xxxxxxxxxxxxxxxxxxxxxxxx")
	default:
		fmt.Println(guess)
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
