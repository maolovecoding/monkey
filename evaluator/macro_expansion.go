package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

// TODO 如何实现嵌套的宏定义呢
// 定义宏
func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}
	for i, statement := range program.Statements {
		if isMacroDefinition(statement) { // 是否是宏定义语句
			addMacro(statement, env)             // 添加
			definitions = append(definitions, i) // 记录宏定义
		}
	}
	for i := len(definitions) - 1; i >= 0; i-- {
		definitionIndex := definitions[i]
		// 删除宏定义的语句
		program.Statements = append(program.Statements[:definitionIndex], program.Statements[definitionIndex+1:]...)
	}
}

// 是宏定义声明
func isMacroDefinition(node ast.Statement) bool {
	letStatement, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}
	_, ok = letStatement.Value.(*ast.MacroLiteral)
	if !ok {
		return false
	}
	return true
}

// 添加宏定义
func addMacro(stmt ast.Statement, env *object.Environment) {
	letStatement, _ := stmt.(*ast.LetStatement)
	macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)
	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Env:        env,
		Body:       macroLiteral.Body,
	}
	// 宏绑定到名称上了 也就是只有声明时的宏定义才是认为是真的宏定义
	env.Set(letStatement.Name.Value, macro)
}
