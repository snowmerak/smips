package process

import (
	"math"
	"math/bits"

	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

func (p *Process) Eval() bool {
	if p.programCount >= uint64(len(p.opcodes)) {
		return false
	}
	p.programCount++
	op := p.opcodes[p.programCount]
	switch opcode.Op(op.Op()) {
	case opcode.OpAdd:
		hi, lo := bits.Add64(p.stack.Get(op.Ra()), p.stack.Get(op.Rb()), 0)
		p.stack.Set(op.Rd(), lo)
		p.stack.SetHi(hi)
	case opcode.OpSub:
		diff, out := bits.Sub64(p.stack.Get(op.Ra()), p.stack.Get(op.Rb()), 0)
		p.stack.Set(op.Rd(), diff)
		p.stack.SetHi(out)
	case opcode.OpMul:
		hi, lo := bits.Mul64(p.stack.Get(op.Ra()), p.stack.Get(op.Rb()))
		p.stack.Set(op.Rd(), lo)
		p.stack.SetHi(hi)
	case opcode.OpDiv:
		div, rem := bits.Div64(0, p.stack.Get(op.Ra()), p.stack.Get(op.Rb()))
		p.stack.Set(op.Rd(), div)
		p.stack.SetHi(rem)
	case opcode.OpMod:
		div, rem := bits.Div64(0, p.stack.Get(op.Ra()), p.stack.Get(op.Rb()))
		p.stack.Set(op.Rd(), rem)
		p.stack.SetHi(div)
	case opcode.OpAnd:
		p.stack.Set(op.Rd(), p.stack.Get(op.Ra())&p.stack.Get(op.Rb()))
	case opcode.OpOr:
		p.stack.Set(op.Rd(), p.stack.Get(op.Ra())|p.stack.Get(op.Rb()))
	case opcode.OpXor:
		p.stack.Set(op.Rd(), p.stack.Get(op.Ra())^p.stack.Get(op.Rb()))
	case opcode.OpNot:
		p.stack.Set(op.Rd(), ^p.stack.Get(op.Ra()))
	case opcode.OpShl:
		p.stack.Set(op.Rd(), p.stack.Get(op.Ra())<<p.stack.Get(op.Rb()))
	case opcode.OpShr:
		p.stack.Set(op.Rd(), p.stack.Get(op.Ra())>>p.stack.Get(op.Rb()))
	case opcode.OpEq:
		if p.stack.Get(op.Ra()) == p.stack.Get(op.Rb()) {
			p.programCount = p.stack.Get(op.Rd())
		}
	case opcode.OpNeq:
		if p.stack.Get(op.Ra()) != p.stack.Get(op.Rb()) {
			p.programCount = p.stack.Get(op.Rd())
		}
	case opcode.OpFAdd:
		p.stack.Set(op.Rd(), math.Float64bits(math.Float64frombits(p.stack.Get(op.Ra()))+math.Float64frombits(p.stack.Get(op.Rb()))))
	case opcode.OpFSub:
		p.stack.Set(op.Rd(), math.Float64bits(math.Float64frombits(p.stack.Get(op.Ra()))-math.Float64frombits(p.stack.Get(op.Rb()))))
	case opcode.OpFMul:
		p.stack.Set(op.Rd(), math.Float64bits(math.Float64frombits(p.stack.Get(op.Ra()))*math.Float64frombits(p.stack.Get(op.Rb()))))
	case opcode.OpFDiv:
		p.stack.Set(op.Rd(), math.Float64bits(math.Float64frombits(p.stack.Get(op.Ra()))/math.Float64frombits(p.stack.Get(op.Rb()))))
	case opcode.OpSet:
		p.stack.Set(op.Rd(), op.WideData())
	case opcode.OpJmp:
		p.programCount = p.stack.Get(op.Ra())
	case opcode.OpAlloc:
		p.stack.Set(op.Rd(), p.memory.Alloc(p.stack.Get(op.Ra())))
	case opcode.OpLoad:
		size := uint64(op.Data() & 0b00000111)
		data := uint64(0)
		for i := uint64(0); i < size; i++ {
			data = data<<8 | uint64(p.memory.Get(p.stack.Get(op.Ra()), p.stack.Get(op.Rb())+i))
		}
		p.stack.Set(op.Rd(), data)
	case opcode.OpStore:
		size := uint64(op.Data() & 0b00000111)
		data := p.stack.Get(op.Rd())
		for i := size - 1; i >= 0; i++ {
			p.memory.Set(p.stack.Get(op.Ra()), p.stack.Get(op.Rb())+i, uint8(data&0xff))
			data = data >> 8
		}
	case opcode.OpLib:
		p.stack.Call(register.New())
		p.libs[op.Ra()].Execute()
		p.stack.Plunder()
	case opcode.OpCall:
		SystemCall(p, op)
	}
	return true
}

func (p *Process) Execute() {
	for p.Eval() {
	}
}
