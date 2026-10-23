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
	op := p.opcodes[p.programCount]
	p.programCount++
	switch opcode.Op(op.Op()) {
	case opcode.OpAdd:
		lo, hi := bits.Add64(p.stack.Get(op.Ra()), p.stack.Get(op.Rb()), 0)
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
	case opcode.OpFMod:
		p.stack.Set(op.Rd(), math.Float64bits(math.Mod(math.Float64frombits(p.stack.Get(op.Ra())), math.Float64frombits(p.stack.Get(op.Rb())))))
	case opcode.OpSet:
		p.stack.Set(op.Rd(), op.WideData())
	case opcode.OpJmp:
		p.programCount = p.stack.Get(op.Ra())
	case opcode.OpAlloc:
		p.stack.Set(op.Rd(), p.memory.Alloc(p.stack.Get(op.Ra())))
	case opcode.OpLoad:
		size := uint64(op.Size())
		data := uint64(0)
		addr, offset := p.stack.Get(op.Ra()), p.stack.Get(op.Rb())
		for i := uint64(0); i < size; i++ {
			data <<= 8
			if i <= math.MaxUint64-offset {
				data |= uint64(p.memory.Get(addr, offset+i))
			}
		}
		p.stack.Set(op.Rd(), data)
	case opcode.OpStore:
		size := uint64(op.Size())
		data := p.stack.Get(op.Rd())
		addr, offset := p.stack.Get(op.Ra()), p.stack.Get(op.Rb())
		length := p.memory.Size(addr)
		if size > 0 && (size > length || offset > length-size) {
			panic("process: store out of range")
		}
		for i := size; i > 0; i-- {
			p.memory.Set(addr, offset+i-1, uint8(data&0xff))
			data = data >> 8
		}
	case opcode.OpLib:
		if int(op.Ra()) >= len(p.libs) {
			panic("process: library index out of range")
		}
		p.call(p.libs[op.Ra()])
	default:
		SystemCall(p, op)
	}
	return true
}

// call runs registered code on a new frame in this runtime. The registered
// Process is a code source; its PC, registers, memory and files are untouched.
func (p *Process) call(lib *Process) {
	opcodes, libs, pc := p.opcodes, p.libs, p.programCount
	p.stack.Call(register.New())
	returned := false
	defer func() {
		p.opcodes, p.libs, p.programCount = opcodes, libs, pc
		if returned {
			p.stack.Plunder()
		} else {
			// A failed child must not publish partial return values.
			p.stack.Pop().Return()
		}
	}()
	p.opcodes, p.libs, p.programCount = lib.opcodes, lib.libs, 0
	p.Execute()
	returned = true
}

func (p *Process) Execute() {
	for p.Eval() {
	}
}
