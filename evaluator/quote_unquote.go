package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

// quote函数 只能有一个参数
func quote(node ast.Node) object.Object {
	return &object.Quote{
		Node: node,
	}
}
