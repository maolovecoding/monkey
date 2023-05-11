package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// 语法解析器对象
type Parser struct {
	l         *lexer.Lexer // 词法解析对象
	curToken  token.Token  // 当前的token 词法单元
	peekToken token.Token  // 偷看的下一个token 如果上一个token的信息不够做决策，需要根据该字段来做决策
	errors    []string     // 错误
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
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
	default:
		return nil
	}
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
