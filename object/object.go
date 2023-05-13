package object

type ObjectType string

// 对象表示
type Object interface {
	Type() ObjectType // 类型
	Inspect() string  // 检查
}
