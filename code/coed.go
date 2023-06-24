package code

import (
	"encoding/binary"
	"fmt"
)

// 定义字节码格式

// 指令
type Instructions []byte

// 操作码类型
type Opcode byte

// 具体的操作码
const (
	OpConstant Opcode = iota // push 使用操作数作为索引检索常量并压栈
)

type Definition struct {
	Name          string // 操作码名称
	OperandWidths []int  // 操作数占用的字节数
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // 操作数2字节宽 最大值 65535
}

// 方便查看指令 操作码 + 操作数
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// 编译字节码 操作码 操作数任意数量
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}
	instructionLen := 1
	// 看指令的长度
	for _, w := range def.OperandWidths {
		instructionLen += w
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}
