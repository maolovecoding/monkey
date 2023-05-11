package token

const (
	ILLEGAL = "ILLEGAL" // 未知的类型 非法类型等
	EOF     = "EOF"     // end of file
	// 标识符 + 字面量
	IDENT = "IDENT" // 变量标识符
	INT   = "INT"   // 数字
	// 运算符
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	// 等于 不等于运算符 为什么在这里指出来？ 因为这是两个字符构成的 也就是双字符的运算符
	EQ     = "=="
	NOT_EQ = "!="
	// TODO 支持更多双字符运算符 >= <= && || 如果支持二进制 还可以支持 >> << 那么应该抽取成方法了 比如 makeTwoCharToken
	// 分割符
	COMMA     = ","
	SEMICOLON = ";"
	// 左右小括号 大括号
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	// 关键字
	LET      = "LET"
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type TokenType string

type Token struct {
	Type    TokenType // 类型
	Literal string    // 字面量值
}

// 定义的关键字获取对应的类型
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// 检查关键字 如果标识符ident是关键字 那就返回对应的常量 否则返回标识符IDENT
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
