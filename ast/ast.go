package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string // 返回值是与其管理的词法单元的字面量 主要是为了方便测试
}

// 声明 语句
type Statement interface {
	Node
	statementNode()
}

// 表达式
type Expression interface {
	Node
	expressionNode()
}

// 程序 也是根
type Program struct {
	Statements []Statement // 切片
}
type LetStatement struct {
	Token token.Token // 词法单元 LET
	Name  *Identifier // 标识符
	Value Expression  // 产生值的表达式
}

// 实现接口
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token // IDENT 词法单元
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
