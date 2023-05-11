package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// 解析函数类型
type (
	prefixParseFn func() ast.Expression               // 前缀解析函数
	infixParseFn  func(ast.Expression) ast.Expression // 中缀解析函数 参数是中缀表达式左侧的内容
	// TODO 支持后缀表达式？
)

// 表达式优先级 从上到下依次递增
const (
	_ int = iota // 空白标识符
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // -X or !X
	CALL        // add()
)

// 优先级表  词法类型关联对应的优先级
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Parser 语法解析器对象
type Parser struct {
	l              *lexer.Lexer                      // 词法解析对象
	curToken       token.Token                       // 当前的token 词法单元
	peekToken      token.Token                       // 偷看的下一个token 如果上一个token的信息不够做决策，需要根据该字段来做决策
	errors         []string                          // 错误
	prefixParseFns map[token.TokenType]prefixParseFn // 词法单元类型关联对应的解析函数
	infixParseFns  map[token.TokenType]infixParseFn
}

// New 创建一个parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// 关联解析函数
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIndentifier)      // 标识符
	p.registerPrefix(token.INT, p.parseIntegerLiteral)     // 数字
	p.registerPrefix(token.BANG, p.parsePrefixExpression)  // / !
	p.registerPrefix(token.MINUS, p.parsePrefixExpression) // -
	// 中缀表达式
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	// 读取两个词法单元 设置 curToken peekToken
	p.nextToken()
	p.nextToken()
	return p
}

// 同时移动词法单元到下一位
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// 开始语法解析 生成ast
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement() // 去处理语句
		if stmt != nil {
			program.Statements = append(program.Statements, stmt) // 追加生成的语句
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement() // 处理let语句
	case token.RETURN:
		return p.parseReturnStatement() // 处理return
	default: // 不是语句 那就开始解析表达式了
		return p.parseExpressionStatement()
	}
}

// parseExpressionStatement 解析表达式语句
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) { // 是分号 跳过 不是分号也无所谓
		p.nextToken()
	}
	return stmt
}

// parseExpression 解析表达式 参数是优先级 是递归的入口 比如前缀表达式处理了 -这种前缀 右操作符还需要处理 又来到这里
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		// 没有对应的前缀表达式解析函数
		p.noPrefixParseFnError(p.curToken.Type) // 收集错误
		return nil
	}
	leftExp := prefix()
	// TODO 调试流程
	// 不是分号词法 且优先级上一个运算符优先级低于当前优先级
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp // 没有右侧词法对应的解析函数
		}
		p.nextToken()
		leftExp = infix(leftExp) // 基于左侧表达式，构建完整中缀表达式
	}
	return leftExp
}

// 解析let 语句
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken} // LET let
	if !p.expectPeek(token.IDENT) {
		return nil // 查看下一个期望词法单元是否是期望的 标识符
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil // 是否是 =词法单元
	}
	// TODO 跳过对求值表达式的处理 直到遇到分号
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	// TODO 跳过对求值表达式的处理 直到遇到分号
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// 当前的tokenTyoe是否是想要的
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 查看下一个tokenType是否是想要的
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 断言函数
// 是否是期待的下一个token
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		// 添加错误 不是期望的tokenType
		p.peekError(t)
		return false
	}
}

// 查看错误
func (p *Parser) Errors() []string {
	return p.errors
}

// 添加错误
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead.", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// registerPrefix 工具方法 注册解析函数
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// parseIndentifier 标识符解析函数
func (p *Parser) parseIndentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// 解析为数字字面量
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	// 字面量是字符串形式的 "10" => 转换为 10
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// 没有前缀表达式对应的解析函数 错误收集
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found.", t)
	p.errors = append(p.errors, msg)
}

// - ! 对应的前缀表达式解析函数
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken() // 向下移动 指向前缀表达式的右操作数
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

// peekPrecedence 下一个词法单元对应的优先级
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p // 获取词法类型对应的优先级
	}
	return LOWEST // 默认最低优先级
}

// curPrecedence 当前词法单元对应的优先级
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p // 获取词法类型对应的优先级
	}
	return LOWEST // 默认最低优先级
}

// parseInfixExpression 解析中缀表达式
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence) // 解析右半部分
	return expression
}
