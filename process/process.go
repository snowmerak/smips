package process

import (
	"github.com/snowmerak/smips/files"
	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

type Process struct {
	files   files.Files
	memory  memory.Memory
	stack   register.Stack
	opcodes []opcode.Opcode
	libs    []*Process

	programCount uint64
	paramCount   uint16
}

func New(opcodes ...opcode.Opcode) *Process {
	p := &Process{
		files:   files.New(),
		memory:  memory.New(),
		stack:   register.NewStack(),
		opcodes: opcodes,
		libs:    make([]*Process, 0),

		programCount: 0,
		paramCount:   0,
	}
	p.stack.Push(register.New())
	return p
}

func (p *Process) AddLib(lib *Process) *Process {
	p.libs = append(p.libs, lib)
	return p
}

func (p *Process) AddStringArgs(args ...String) *Process {
	for _, arg := range args {
		addr := WriteString(&p.memory, arg)
		p.stack.SetParams(p.paramCount, addr)
	}
	return p
}
