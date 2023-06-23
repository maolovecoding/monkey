package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

// quote函数 只能有一个参数 不对参数求值
func quote(node ast.Node) object.Object {
	return &object.Quote{
		Node: node,
	}
}

// unquote 和quote相反 可以对调用的参数求值 但是该函数只能在quote函数中使用
// func unquote() {}
