package ast

// 修改ast节点
type ModifierFunc func(Node) Node

// 遍历的节点 修改节点的函数
func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	// 有子节点可以递归修改
	case *Program: // 根
		for i, statemenrt := range node.Statements {
			// 修改原AST的节点
			node.Statements[i], _ = Modify(statemenrt, modifier).(Statement)
		}
	case *ExpressionStatement: // 表达式语句
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)
	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}
	case *BlockStatement:
		for i := range node.Statements {
			node.Statements[i], _ = Modify(node.Statements[i], modifier).(Statement)
		}
	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
	case *LetStatement:
		node.Value = Modify(node.Value, modifier).(Expression)
	case *FunctionLiteral:
		for i := range node.Parameters {
			node.Parameters[i] = Modify(node.Parameters[i], modifier).(*Identifier) // 处理函数字面量参数
		}
		node.Body = Modify(node.Body, modifier).(*BlockStatement) // 函数体的处理
	case *ArrayLiteral:
		for i := range node.Elements {
			node.Elements[i] = Modify(node.Elements[i], modifier).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, value := range node.Pairs {
			newKey := Modify(key, modifier).(Expression)
			newValue := Modify(value, modifier).(Expression)
			newPairs[newKey] = newValue
		}
		node.Pairs = newPairs
	}
	// 没有子节点了 直接修改
	return modifier(node)
}
