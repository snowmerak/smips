package process

import (
	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

func executeI(pc uint64, code *opcode.OpCode, memory *memory.Memory, register *[]*register.Register) {
	i := opcode.I(*code)
	switch i.Function() {
	}
}
