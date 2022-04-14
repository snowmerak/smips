package process

import (
	"math/bits"

	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

func executeR(pc uint64, code *opcode.OpCode, memory *memory.Memory, registers *[]*register.Register) {
	r := opcode.R(*code)
	reg := (*registers)[len(*registers)-1]
	a := reg.Get(r.RegisterSourceA())
	b := reg.Get(r.RegisterSourceB())
	switch r.Function() {
	case opcode.RAdd:
		sum, carry := bits.Add64(a, b, 0)
		reg.Set(r.RegisterDestination(), sum)
		reg.SetHi(carry)
	case opcode.RSub:
		sum, carry := bits.Sub64(a, b, 0)
		reg.Set(r.RegisterDestination(), sum)
		reg.SetHi(carry)
	case opcode.RMul:
		hi, lo := bits.Mul64(a, b)
		reg.Set(r.RegisterDestination(), lo)
		reg.SetHi(hi)
	case opcode.RDiv:
		div, rem := bits.Div64(0, a, b)
		reg.Set(r.RegisterDestination(), div)
		reg.SetHi(rem)
	case opcode.RRem:
		rem := bits.Rem64(0, a, b)
		reg.Set(r.RegisterDestination(), rem)
		reg.SetHi(rem)
	case opcode.RAnd:
		and := a & b
		reg.Set(r.RegisterDestination(), and)
	case opcode.ROr:
		or := a | b
		reg.Set(r.RegisterDestination(), or)
	case opcode.RXor:
		xor := a ^ b
		reg.Set(r.RegisterDestination(), xor)
	case opcode.RNor:
		nor := ^(a | b)
		reg.Set(r.RegisterDestination(), nor)
	case opcode.RShiftLeft:
		shl := a << b
		reg.Set(r.RegisterDestination(), shl)
	case opcode.RShiftRightLogical:
		shr := a >> b
		reg.Set(r.RegisterDestination(), shr)
	case opcode.RShiftRightArithmetic:
		sra := int64(a) >> int64(b)
		reg.Set(r.RegisterDestination(), uint64(sra))
	}
}
