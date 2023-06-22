package evaluator

import "monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array: // 数组
				return &object.Integer{
					Value: int64(len(arg.Elements)),
				}
			case *object.String: // 字符串长度
				return &object.Integer{
					Value: int64(len(arg.Value)),
				}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { // 参数长度校验
				return newError("wrong number of arguments.got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ { // 不是数组 不支持first函数
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { // 参数长度校验
				return newError("wrong number of arguments.got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ { // 不是数组 不支持first函数
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rust": { // 返回除去第一个元素的新数组 不会修改原数组
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { // 参数长度校验
				return newError("wrong number of arguments.got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ { // 不是数组 不支持first函数
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{
					Elements: newElements,
				}
			}
			return NULL
		},
	},
	"push": { // 向数组追加元素 返回是追加元素后的数组 原数组是不变的 数组具有不可变性
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 { // 参数长度校验
				return newError("wrong number of arguments.got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ { // 不是数组 不支持first函数
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{
				Elements: newElements,
			}
		},
	},
	"pop": { // 向数组弹出最后一个元素 返回是弹出元素后的数组 原数组是不变的 数组具有不可变性
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { // 参数长度校验
				return newError("wrong number of arguments.got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ { // 不是数组 不支持first函数
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length >= 1 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[0:length-1])
				return &object.Array{
					Elements: newElements,
				}
			}
			return &object.Array{} // 空数组没有元素 返回值还是空数组就行
		},
	},
}
