package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current() // 当前的用户
	if err != nil {
		return
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username) // 获取用户名
	repl.Start(os.Stdin, os.Stdout)
}
