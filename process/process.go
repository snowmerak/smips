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

// AddLib registers code for an internal frame call. Its execution uses the
// caller's memory and files, not the registered Process's runtime state.
func (p *Process) AddLib(lib *Process) *Process {
	if lib == nil {
		panic("process: nil library")
	}
	p.libs = append(p.libs, lib)
	return p
}

// AddArgs appends raw register values to the current frame's parameters.
func (p *Process) AddArgs(args ...uint64) *Process {
	p.checkArgs(len(args))
	for _, arg := range args {
		p.stack.SetParams(p.paramCount, arg)
		p.paramCount++
	}
	return p
}

// GetReturn reads a result from the current frame, including the root after Execute.
func (p *Process) GetReturn(index uint16) uint64 {
	if index >= register.BankSize {
		panic("process: return index out of range")
	}
	return p.stack.GetReturns(index)
}

func (p *Process) AddStringArgs(args ...String) *Process {
	p.checkArgs(len(args))
	for _, arg := range args {
		addr := WriteString(&p.memory, arg)
		p.stack.SetParams(p.paramCount, addr)
		p.paramCount++
	}
	return p
}

func (p *Process) checkArgs(count int) {
	if count > register.BankSize-int(p.paramCount) {
		panic("process: too many arguments")
	}
}
