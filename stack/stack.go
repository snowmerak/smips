package stack

import (
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

type Stack struct {
	opcodes []opcode.OpCode
	reg     register.Register
}
