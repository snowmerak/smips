package process

import (
	"math/bits"

	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

func executeI(pc uint64, code *opcode.OpCode, memory *memory.Memory, register *[]*register.Register) {
	i := opcode.I(*code)
	src := i.RegisterSource()
	dest := i.RegisterDestination()
	switch i.Function() {
	case opcode.ISet:
		(*register)[len(*register)-1].SetV(dest, i.Immediate())
	case opcode.IAdd:
		v := (*register)[len(*register)-1].GetV(src)
		sum, carray := bits.Add64(v, i.Immediate(), 0)
		(*register)[len(*register)-1].SetV(dest, sum)
		(*register)[len(*register)-1].SetHi(carray)
	case opcode.ISub:
		v := (*register)[len(*register)-1].GetV(src)
		diff, carray := bits.Sub64(v, i.Immediate(), 0)
		(*register)[len(*register)-1].SetV(dest, diff)
		(*register)[len(*register)-1].SetHi(carray)
	case opcode.IMul:
		v := (*register)[len(*register)-1].GetV(src)
		hi, lo := bits.Mul64(v, i.Immediate())
		(*register)[len(*register)-1].SetV(dest, lo)
		(*register)[len(*register)-1].SetHi(hi)
	case opcode.IDiv:
		v := (*register)[len(*register)-1].GetV(src)
		div, rem := bits.Div64(0, v, i.Immediate())
		(*register)[len(*register)-1].SetV(dest, div)
		(*register)[len(*register)-1].SetHi(rem)
	case opcode.IRem:
		v := (*register)[len(*register)-1].GetV(src)
		rem := bits.Rem64(0, v, i.Immediate())
		(*register)[len(*register)-1].SetV(dest, rem)
		(*register)[len(*register)-1].SetHi(rem)
	case opcode.ILoadData:
		(*register)[len(*register)-1].SetV(dest, memory.Load((*register)[len(*register)-1].Get(src)))
	case opcode.IStoreData:
		memory.Store((*register)[len(*register)-1].Get(dest), (*register)[len(*register)-1].Get(src))
	}
}
