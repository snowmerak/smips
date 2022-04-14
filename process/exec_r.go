package process

import (
	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

func executeR(pc uint64, code *opcode.OpCode, memory *memory.Memory, register *[]*register.Register) {
	r := opcode.R(*code)
	switch r.Function() {

	}
}
