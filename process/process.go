package process

import (
	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

type Process struct {
	memory        *memory.Memory
	opcodes       []opcode.OpCode
	registerStack []*register.Register
}
