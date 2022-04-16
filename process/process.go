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

// New creates a new process.
func New(opcodes ...opcode.OpCode) *Process {
	return &Process{
		memory:        memory.New(1024),
		opcodes:       opcodes,
		registerStack: []*register.Register{register.New()},
	}
}

// Execute executes the process.
func (p *Process) Execute() {
	for p.pc < uint64(len(p.opcodes)) {
		opcode := p.opcodes[p.pc]
		Execute(p.pc, &opcode, p.memory, &p.registerStack)
		p.pc++
	}
}
