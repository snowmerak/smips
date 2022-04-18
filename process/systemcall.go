package process

import "github.com/snowmerak/smips/opcode"

func SystemCall(p *Process, op opcode.Opcode) {
	switch opcode.Op(op.Op()) {
	case opcode.OpenFile:
		path := ReadString(&p.memory, p.stack.Get(op.Ra()), 0)
		fd := p.files.OpenFile(string(path), uint8(p.stack.Get(op.Rb())))
		p.stack.Set(op.Rd(), uint64(fd))
	case opcode.CloseFile:
		p.stack.Set(op.Rd(), p.files.CloseFile(p.stack.Get(op.Ra())))
	case opcode.ReadFile:
		p.stack.Set(op.Rd(), p.files.ReadFile(p.stack.Get(op.Ra()), op.Data()))
	case opcode.ReadFileAt:
		p.stack.Set(op.Rd(), p.files.ReadFileAt(p.stack.Get(op.Ra()), p.stack.Get(op.Rb()), op.Data()))
	case opcode.WriteFile:
		p.stack.SetHi(p.files.WriteFile(p.stack.Get(op.Ra()), op.Data(), p.stack.Get(op.Rd())))
	case opcode.WriteFileAt:
		p.stack.SetHi(p.files.WriteFileAt(p.stack.Get(op.Ra()), p.stack.Get(op.Rb()), op.Data(), p.stack.Get(op.Rd())))
	}
}
