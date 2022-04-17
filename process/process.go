package process

import (
	"github.com/snowmerak/smips/files"
	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

type Process struct {
	files        files.Files
	memory       memory.Memory
	stack        register.Stack
	programCount uint64
	opcodes      []opcode.Opcode
	libs         []*Process
}

func New(opcodes ...opcode.Opcode) *Process {
	p := &Process{
		files:        files.New(),
		memory:       memory.New(),
		stack:        register.NewStack(),
		programCount: 0,
		opcodes:      opcodes,
		libs:         make([]*Process, 0),
	}
	p.stack.Push(register.New())
	return p
}

func (p *Process) AddLib(lib *Process) *Process {
	p.libs = append(p.libs, lib)
	return p
}
