package token

const (
	ILLEGAL = "ILLEGAL" // 未知的类型 非法类型等
	EOF     = "EOF"     // end of file
	// 标识符 + 字面量
	IDENT = "IDENT" // 变量标识符
	INT   = "INT"   // 数字
	// 运算符
	ASSIGN = "="
	PLUS   = "+"
	// 分割符
	COMMA     = ","
	SEMICOLON = ";"
	// 左右小括号 大括号
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
)

type TokenType string

type Token struct {
	Type    TokenType // 类型
	Literal string    // 字面量值
}
