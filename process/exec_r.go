package process

import (
	"math"
	"math/bits"

	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

func executeR(pc *uint64, code *opcode.OpCode, _ *memory.Memory, registers *[]*register.Register) {
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
	case opcode.RBranchEqual:
		if a == b {
			*pc = uint64((*registers)[len(*registers)-1].Get(r.RegisterDestination()))
		}
	case opcode.RBranchNotEqual:
		if a != b {
			*pc = uint64((*registers)[len(*registers)-1].Get(r.RegisterDestination()))
		}
	case opcode.RBranchLessThan:
		if a < b {
			*pc = uint64((*registers)[len(*registers)-1].Get(r.RegisterDestination()))
		}
	case opcode.RBranchLessThanOrEqual:
		if a <= b {
			*pc = uint64((*registers)[len(*registers)-1].Get(r.RegisterDestination()))
		}
	case opcode.RBranchGreaterThan:
		if a > b {
			*pc = uint64((*registers)[len(*registers)-1].Get(r.RegisterDestination()))
		}
	case opcode.RBranchGreaterThanOrEqual:
		if a >= b {
			*pc = uint64((*registers)[len(*registers)-1].Get(r.RegisterDestination()))
		}
	case opcode.RFloatAdd:
		fa := math.Float64frombits(a)
		fb := math.Float64frombits(b)
		fadd := fa + fb
		reg.Set(r.RegisterDestination(), math.Float64bits(fadd))
	case opcode.RFloatSub:
		fa := math.Float64frombits(a)
		fb := math.Float64frombits(b)
		fsub := fa - fb
		reg.Set(r.RegisterDestination(), math.Float64bits(fsub))
	case opcode.RFloatMul:
		fa := math.Float64frombits(a)
		fb := math.Float64frombits(b)
		fmul := fa * fb
		reg.Set(r.RegisterDestination(), math.Float64bits(fmul))
	case opcode.RFloatDiv:
		fa := math.Float64frombits(a)
		fb := math.Float64frombits(b)
		fdiv := fa / fb
		reg.Set(r.RegisterDestination(), math.Float64bits(fdiv))
	}
}
