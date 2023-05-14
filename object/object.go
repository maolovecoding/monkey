package object

import "fmt"

type ObjectType string // 对象类型

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
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

// 布尔值
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

// Null
type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

// 返回值对象 记录是否返回了值 return 有值就停止求值
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

func (n *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}
