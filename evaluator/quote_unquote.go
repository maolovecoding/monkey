package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/token"
)

// quote函数 只能有一个参数 不对参数求值
func quote(node ast.Node, env *object.Environment) object.Object {
	// 遇到需要处理的node
	node = evalUnquoteCalls(node, env)
	return &object.Quote{
		Node: node,
	}
}

// 处理ast 需要进行求值的ast就会求值
func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node // 不是uniquote调用
		}
		// uniquote调用 需要对参数求值的
		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node // 不是调用表达式
		}
		if len(call.Arguments) != 1 {
			return node // 参数不是1
		}
		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToASTNode(unquoted)
	})
}

// 转换 object为node
func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{
			Token: t,
			Value: obj.Value,
		}
	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{
				Type:    token.TRUE,
				Literal: "true",
			}
		} else {
			t = token.Token{
				Type:    token.FALSE,
				Literal: "false",
			}
		}
		return &ast.Boolean{
			Token: t,
			Value: obj.Value,
		}
	case *object.Quote:
		return obj.Node // unique(quote(node)) 嵌套quote 那么被嵌套的quote已经处理过node 这里不要二次处理！
		// TODO 错误处理？
	default:
		return nil
	}
}

// 是否是unquote函数调用
func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExpression.Function.TokenLiteral() == "unquote"
}

// unquote 和quote相反 可以对调用的参数求值 但是该函数只能在quote函数中使用
// func unquote() {}
