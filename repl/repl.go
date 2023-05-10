// read-eval-print Loop 读取求值打印循环
// 我们可以从控制台获取到输入的内容，也就是所谓的源代码，将读取的内容传递给词法分析器 然后输出词法分析器生成的词法单元，直到遇到EOF
package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> " // prompt
// 读取命令行输入的源代码
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok) // 输出到控制台 以go风格输出
		}
	}
}
