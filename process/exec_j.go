package process

import (
	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

func Execute(pc uint64, code *opcode.OpCode, memory *memory.Memory, registers *[]*register.Register) {
	switch code.Type() {
	case opcode.TypeR:
		executeR(pc, code, memory, registers)
	case opcode.TypeI:
		executeI(pc, code, memory, registers)
	case opcode.TypeJ:
		executeJ(pc, code, memory, registers)
	}
}

func executeJ(pc uint64, code *opcode.OpCode, _ *memory.Memory, registers *[]*register.Register) {
	j := opcode.J(*code)
	switch j.Function() {
	case opcode.Jump:
		pc = j.Address()
	case opcode.JumpAndAddRegister:
		*registers = append(*registers, register.New())
		pc = j.Address()
	case opcode.JumpAndRemoveRegister:
		*registers = (*registers)[:len(*registers)-1]
		pc = j.Address()
	case 3:
	}
}
