package process

import (
	"math/bits"

	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

func executeI(pc *uint64, code *opcode.OpCode, memory *memory.Memory, registers *[]*register.Register) {
	i := opcode.I(*code)
	src := i.RegisterSource()
	dest := i.RegisterDestination()
	switch i.Function() {
	case opcode.ISet:
		(*registers)[len(*registers)-1].SetV(dest, i.Immediate())
	case opcode.IAdd:
		v := (*registers)[len(*registers)-1].GetV(src)
		sum, carray := bits.Add64(v, i.Immediate(), 0)
		(*registers)[len(*registers)-1].SetV(dest, sum)
		(*registers)[len(*registers)-1].SetHi(carray)
	case opcode.ISub:
		v := (*registers)[len(*registers)-1].GetV(src)
		diff, carray := bits.Sub64(v, i.Immediate(), 0)
		(*registers)[len(*registers)-1].SetV(dest, diff)
		(*registers)[len(*registers)-1].SetHi(carray)
	case opcode.IMul:
		v := (*registers)[len(*registers)-1].GetV(src)
		hi, lo := bits.Mul64(v, i.Immediate())
		(*registers)[len(*registers)-1].SetV(dest, lo)
		(*registers)[len(*registers)-1].SetHi(hi)
	case opcode.IDiv:
		v := (*registers)[len(*registers)-1].GetV(src)
		div, rem := bits.Div64(0, v, i.Immediate())
		(*registers)[len(*registers)-1].SetV(dest, div)
		(*registers)[len(*registers)-1].SetHi(rem)
	case opcode.IRem:
		v := (*registers)[len(*registers)-1].GetV(src)
		rem := bits.Rem64(0, v, i.Immediate())
		(*registers)[len(*registers)-1].SetV(dest, rem)
		(*registers)[len(*registers)-1].SetHi(rem)
	case opcode.ILoadData:
		(*registers)[len(*registers)-1].SetV(dest, memory.Load((*registers)[len(*registers)-1].Get(src)))
	case opcode.IStoreData:
		memory.Store((*registers)[len(*registers)-1].Get(dest), (*registers)[len(*registers)-1].Get(src))
	case opcode.IJump:
		*pc = (*registers)[len(*registers)-1].Get(dest)
	}
}
