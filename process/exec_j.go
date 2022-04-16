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
	case opcode.JumpAndPassParameter:
		newReg := register.New()
		for i := byte(0); i < 32; i++ {
			newReg.SetA(i, (*registers)[len(*registers)-1].GetA(i))
		}
		*registers = append(*registers, newReg)
		pc = j.Address()
	case opcode.JumpAndRemoveRegister:
		last := (*registers)[len(*registers)-1]
		*registers = (*registers)[:len(*registers)-1]
		last.Return()
		pc = j.Address()
	case opcode.JumpAndReturnValue:
		last := (*registers)[len(*registers)-1]
		*registers = (*registers)[:len(*registers)-1]
		for i := byte(0); i < 32; i++ {
			(*registers)[len(*registers)-1].SetV(i, last.GetV(i))
		}
		last.Return()
		pc = j.Address()
	}
}
