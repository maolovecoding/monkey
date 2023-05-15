package lexer

import (
	"monkey/token"
)

// TODO 1. 支持utf-8 2. 支持如表情字符？ ch应该使用rune了

// 词法解析
type Lexer struct {
	input        string // 源代码
	position     int    // 当前正在查看的位置 指向当前查看字符的下标
	readPosition int    // 当前字符的下一个字符所在的位置 下标
	ch           byte   // 当前正在查看的字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// =================== 实现的方法 ===============
// 读取下一个字符 readPosition指向下一个字符的位置了
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 读取到末尾
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 获取下一个token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	// 跳过空白字符
	l.skipWhitespace()
	switch l.ch {
	case '=':
		// 多看一个字符 是否可以组成 ==
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch) // = + =
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch) // / ! + =
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"': // string
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = token.EOF
		tok.Literal = "" // 文件末尾了
	default: // 不是可直接识别的字符 检查是否是标识符
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal) // let标识符拿到其关键字LET类型
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch) // 未知的字符
		}
	}
	l.readChar()
	return tok
}

// 读取标识符 且后移
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) { // 只要是字符一直读取 start:end
		l.readChar()
	}
	return l.input[position:l.position] // position所在的字符不是letter了
}

// 跳过空白字符 也可以说"消费 吃掉"
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// 读取一个数字
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) { // 只要是字符一直读取 start:end
		l.readChar()
	}
	return l.input[position:l.position] // position所在的字符不是letter了
}

// 偷看当前字符的下一个字符是什么 并返回
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// TODO 支持 \n 转义字符 支持没有结束 " 报错 ？
// 拿到一个字符串 "add" => add
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// ============= 函数 ================

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9' // TODO 如何支持浮点型？ 二进制八进制十六进制
}
