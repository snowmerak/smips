package process

import (
	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

// Process is a process. lol.
type Process struct {
	memory        *memory.Memory
	opcodes       []opcode.OpCode
	registerStack []*register.Register

	pc uint64
}
