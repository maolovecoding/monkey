package ast

import (
	"bytes"
	"monkey/token"
	"strings"
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

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
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

// 中缀表达式
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String() // 5+5 => (5 + 5)
}

// 布尔类型
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}

// 块级语句
type BlockStatement struct {
	Token      token.Token // "{" 词法单元
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// if表达式
type IfExpression struct {
	Token       token.Token     // "if" 词法单元
	Condition   Expression      // 判断条件
	Consequence *BlockStatement // 结果 根据条件分别指向两种结果
	Alternative *BlockStatement // 可替代的结果 else语句的结果
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		// 有else语句
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// 函数字面量
type FunctionLiteral struct {
	Token      token.Token     // fn 词法单元
	Parameters []*Identifier   // 参数列表
	Body       *BlockStatement // 函数体
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())
	return out.String()
}

// 调用表达式
type CallExpression struct {
	Token     token.Token  // "(" 词法单元
	Function  Expression   // 标识符或者函数字面量 add() => add fn(){}() => fn(){}
	Arguments []Expression // 调用参数
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// 字符串字面量
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}
