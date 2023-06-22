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
	INDEX       // arr[index] 索引运算符 优先级最高
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
	token.LPAREN:   CALL,  // 函数调用表达式 具备最高优先级
	token.LBRACKET: INDEX, // [ 索引表达式访问优先级
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
	p.registerPrefix(token.IDENT, p.parseIndentifier)        // 标识符
	p.registerPrefix(token.INT, p.parseIntegerLiteral)       // 数字
	p.registerPrefix(token.BANG, p.parsePrefixExpression)    // / !
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)   // -
	p.registerPrefix(token.TRUE, p.parseBoolean)             // true
	p.registerPrefix(token.FALSE, p.parseBoolean)            // false
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression) // (
	p.registerPrefix(token.IF, p.parseIfExpression)          // if
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral) // fn
	p.registerPrefix(token.STRING, p.parseStringLiteral)     // string
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)    // array [
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
	p.registerInfix(token.LPAREN, p.parseCallExpression)    // "(" 解析函数
	p.registerInfix(token.LBRACKET, p.parseIndexExpression) // [ 解析函数
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
	p.nextToken() // = 跳过
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
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

// 解析布尔字面量
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE), // true or false
	}
}

// 解析 ( 左括号
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil // 不是右括号 没办法和左括号配对 说明代码有问题
	}
	return exp
}

// 解析 if
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LPAREN) {
		// 不是小括号 语法错误 if ()
		return nil
	}
	p.nextToken()                                    // 跳过 小括号
	expression.Condition = p.parseExpression(LOWEST) // 条件
	if !p.expectPeek(token.RPAREN) {
		return nil // 右侧小括号
	}
	if !p.expectPeek(token.LBRACE) { // 左侧大括号
		return nil
	}
	expression.Consequence = p.parseBlockStatement() // 结果 if
	if p.peekTokenIs(token.ELSE) {
		// 有没有 else
		p.nextToken() // 跳过 }
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

// 解析块级语句
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}
	p.nextToken() // 跳过 {
	// 不是 } 不是结束符
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// 解析函数字面量
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LPAREN) { // (
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) { // {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

// 解析函数参数
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		// 没有参数
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	identifiers = append(identifiers, ident) // 第一个参数
	for p.peekTokenIs(token.COMMA) {         // 下一个是逗号 后面还有参数
		p.nextToken() // 越过当前参数和逗号
		p.nextToken()
		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil // 没有右小括号
	}
	return identifiers
}

// 解析函数调用 给 ( 注册解析函数
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	exp.Arguments = p.parseExpressionList(token.RPAREN) // 结束末尾是 )
	return exp
}

// 函数调用参数 函数废弃 使用 parseExpressionList
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(token.RPAREN) { // ")" callFn() 没有参数
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST)) // 解析args
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) { // ")"
		return nil
	}
	return args
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// parseArrayLiteral 解析数组字面量
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{
		Token: p.curToken,
	}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

// 在 parseCallArguments 的基础上修改的通用版本 解析表达式列表 拿到表达式列表 函数调用参数也在这里了
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) { // 到数组末尾
		p.nextToken()
		return list
	}
	p.nextToken() // 跳过逗号
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil // 语法错误
	}
	return list
}

// 解析索引表达式 把 [ 当做中缀运算符 arr就是左操作数 0 就是右操作数 arr[0]
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}
