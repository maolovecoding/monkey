package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

// true和false创建引用 只有两个实例
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// 语句
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
		// 表达式
	// 整数字面量
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
	}
	return nil
}

// 解析源代码 根ast
func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)
		// 是return语句 不在继续执行后续的求值
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value // 返回结果
		}
	}
	return result
}

// 解析块级语句 if () {xx}
func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result // 不解包return 会在evalProgram想解包拿到return的value 嵌套块的解决
		}
	}
	return result
}

// 返回缓存的布尔值对象
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// 前缀表达式求值
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default: // 不支持该运算符
		return NULL
	}
}

// 返回 表达式的值取反的真值 布尔值取反
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// 返回 表达式的值的相反数
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ { // 不是整数
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{
		Value: -value,
	}
}

// 中缀表达式
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==": // 对于布尔类型的 == != 比较 因为底层的true和false是相同的object表示 这样是很合适比较的
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}

// 都是整数的中缀表达式
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return NULL
	}
}

// if-else语句
func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	}
	if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}
	// 没有else语句
	return NULL
}

// 是否是真值 默认返回true
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
