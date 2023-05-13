package object

import "fmt"

type ObjectType string // 对象类型

const (
	INTEGER_OBJ = "INTEGER"
)

// 对象表示
type Object interface {
	Type() ObjectType // 类型
	Inspect() string  // 检查
}

// 整数
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
