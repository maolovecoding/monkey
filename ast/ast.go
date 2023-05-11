package ast

import (
	"bytes"
	"monkey/token"
)

type Node interface {
	TokenLiteral() string // 返回值是与其管理的词法单元的字面量 主要是为了方便测试
	String() string       // 方便打印和比较AST节点
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

func (p *Program) String() string {
	var out bytes.Buffer // 缓冲区
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type Identifier struct {
	Token token.Token // IDENT 词法单元
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

// return语句
type ReturnStatement struct {
	Token       token.Token // return 词法单元
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement 表达式语句 不是真正的语句 是仅由表达式构成的语句 一层封装
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// 数字字面量
type IntegerLiteral struct {
	Token token.Token
	Value int64 // 记录数字字面量了 token.INT 关联解析函数 parseIntegerLiteral
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// 前缀表达式
type PrefixExpression struct {
	Token    token.Token // 前缀词法单元 ! -
	Operator string      // 操作符 ! -
	Right    Expression  // 右表达式
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String() // -5 => (-5)
}
