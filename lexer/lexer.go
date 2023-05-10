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
func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.ch = 0 // 读取到末尾
	} else {
		lexer.ch = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

// 获取下一个token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
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
	case 0:
		tok.Type = token.EOF
		tok.Literal = "" // 文件末尾了
	}
	l.readChar()
	return tok
}

// ============= 函数 ================

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
