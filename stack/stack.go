package stack

import (
	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

type Stack struct {
	globalMemory *memory.Memory
	opcodes      []opcode.OpCode
	reg          register.Register
}
