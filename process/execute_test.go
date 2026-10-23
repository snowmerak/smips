package process

import (
	"math"
	"testing"

	"github.com/snowmerak/smips/opcode"
	"github.com/snowmerak/smips/register"
)

func regular(op opcode.Op, data uint8, ra, rb, rd uint16) opcode.Opcode {
	return opcode.NewRegularOpCode(data, ra, rb, rd, uint8(op))
}

func immediate(value uint64, rd uint16) opcode.Opcode {
	return opcode.NewWideOpCode(value, rd, uint8(opcode.OpSet))
}

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	fn()
}

func TestEvalAndTermination(t *testing.T) {
	p := New(immediate(42, register.ReturnsBase))
	if !p.Eval() || p.GetReturn(0) != 42 || p.Eval() || p.Eval() {
		t.Fatal("first/last instruction or completed evaluation is incorrect")
	}
	if New().Eval() {
		t.Fatal("empty program executed an instruction")
	}
}

func TestIntegerOperations(t *testing.T) {
	cases := []struct {
		name             string
		op               opcode.Op
		a, b, result, hi uint64
	}{
		{"add", opcode.OpAdd, 40, 2, 42, 0},
		{"carry", opcode.OpAdd, math.MaxUint64, 1, 0, 1},
		{"subtract", opcode.OpSub, 42, 2, 40, 0},
		{"borrow", opcode.OpSub, 0, 1, math.MaxUint64, 1},
		{"multiply", opcode.OpMul, math.MaxUint64, 2, math.MaxUint64 - 1, 1},
		{"divide", opcode.OpDiv, 17, 5, 3, 2},
		{"modulo", opcode.OpMod, 17, 5, 2, 3},
		{"and", opcode.OpAnd, 0xc, 0xa, 8, 0},
		{"or", opcode.OpOr, 0xc, 0xa, 0xe, 0},
		{"xor", opcode.OpXor, 0xc, 0xa, 6, 0},
		{"not", opcode.OpNot, 0, 0, math.MaxUint64, 0},
		{"left shift", opcode.OpShl, 3, 4, 48, 0},
		{"right shift", opcode.OpShr, 48, 4, 3, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := New(regular(tc.op, 0, register.ParamsBase, register.ParamsBase+1, register.ReturnsBase)).AddArgs(tc.a, tc.b)
			p.Execute()
			if p.GetReturn(0) != tc.result || p.stack.GetHi() != tc.hi {
				t.Fatalf("result/Hi = %d/%d, want %d/%d", p.GetReturn(0), p.stack.GetHi(), tc.result, tc.hi)
			}
		})
	}
}

func TestFloatOperations(t *testing.T) {
	for _, tc := range []struct {
		op   opcode.Op
		want float64
	}{
		{opcode.OpFAdd, 7.5}, {opcode.OpFSub, 3.5},
		{opcode.OpFMul, 11}, {opcode.OpFDiv, 2.75}, {opcode.OpFMod, 1.5},
	} {
		p := New(regular(tc.op, 0, register.ParamsBase, register.ParamsBase+1, register.ReturnsBase)).
			AddArgs(math.Float64bits(5.5), math.Float64bits(2))
		p.Execute()
		if got := math.Float64frombits(p.GetReturn(0)); got != tc.want {
			t.Fatalf("operation %d = %g, want %g", tc.op, got, tc.want)
		}
	}
}

func TestBranches(t *testing.T) {
	for _, tc := range []struct {
		name  string
		op    opcode.Op
		a, b  uint64
		taken bool
	}{
		{"equal taken", opcode.OpEq, 7, 7, true},
		{"equal fallthrough", opcode.OpEq, 7, 8, false},
		{"unequal taken", opcode.OpNeq, 7, 8, true},
		{"unequal fallthrough", opcode.OpNeq, 7, 7, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := New(
				regular(tc.op, 0, register.ParamsBase, register.ParamsBase+1, register.ParamsBase+2),
				immediate(99, register.ReturnsBase),
				immediate(42, register.ReturnsBase+1),
			).AddArgs(tc.a, tc.b, 2)
			p.Execute()
			if (p.GetReturn(0) == 0) != tc.taken || p.GetReturn(1) != 42 {
				t.Fatalf("branch results = %d, %d", p.GetReturn(0), p.GetReturn(1))
			}
		})
	}
}

func TestBackwardJump(t *testing.T) {
	p := New(
		regular(opcode.OpSub, 0, register.ParamsBase, 1, register.ParamsBase),
		regular(opcode.OpNeq, 0, register.ParamsBase, 0, 0),
		immediate(9, register.ReturnsBase),
	).AddArgs(3)
	for i := 0; i < 10; i++ {
		if !p.Eval() {
			if p.GetReturn(0) != 9 || p.stack.GetParams(0) != 0 {
				t.Fatal("backward branch did not execute the loop")
			}
			return
		}
	}
	t.Fatal("loop did not terminate within 10 steps")
}

func TestJump(t *testing.T) {
	p := New(
		regular(opcode.OpJmp, 0, register.ParamsBase, 0, 0),
		immediate(99, register.ReturnsBase),
		immediate(42, register.ReturnsBase+1),
	).AddArgs(2)
	p.Execute()
	if p.GetReturn(0) != 0 || p.GetReturn(1) != 42 {
		t.Fatal("jump did not target the exact instruction index")
	}
}

func TestMemoryInstructions(t *testing.T) {
	for _, size := range []uint8{0, 1, 7, 8} {
		const value = uint64(0x0123456789abcdef)
		p := New(
			regular(opcode.OpAlloc, 0, register.ParamsBase, 0, register.TemporaryBase),
			regular(opcode.OpStore, size, register.TemporaryBase, register.ParamsBase+1, register.ParamsBase+2),
			regular(opcode.OpLoad, size, register.TemporaryBase, register.ParamsBase+1, register.ReturnsBase),
		).AddArgs(10, 2, value)
		p.Execute()
		var want uint64
		if size == 8 {
			want = value
		} else if size > 0 {
			want = value & (uint64(1)<<(8*size) - 1)
		}
		if p.GetReturn(0) != want {
			t.Fatalf("size %d: load = %x, want %x", size, p.GetReturn(0), want)
		}
		for i := uint8(0); i < size; i++ {
			expected := uint8(value >> (8 * (size - i - 1)))
			if p.memory.Get(0, 2+uint64(i)) != expected {
				t.Fatalf("size %d: memory is not big-endian", size)
			}
		}
	}
}

func TestInvalidStoreDoesNotPartiallyWrite(t *testing.T) {
	for _, offset := range []uint64{1, math.MaxUint64} {
		p := New(regular(opcode.OpStore, 2, register.ParamsBase, register.ParamsBase+1, register.ParamsBase+2)).AddArgs(0, offset, 0xffff)
		p.memory.Alloc(2)
		mustPanic(t, func() { p.Eval() })
		if p.memory.Get(0, 0) != 0 || p.memory.Get(0, 1) != 0 {
			t.Fatal("rejected store modified memory")
		}
	}
}

func TestLoadOffsetDoesNotWrap(t *testing.T) {
	p := New(regular(opcode.OpLoad, 2, register.ParamsBase, register.ParamsBase+1, register.ReturnsBase)).AddArgs(0, math.MaxUint64)
	p.memory.Alloc(1)
	p.memory.Set(0, 0, 99)
	p.Execute()
	if p.GetReturn(0) != 0 {
		t.Fatal("out-of-range offset wrapped to the start of the block")
	}
}

func TestInvalidInstructions(t *testing.T) {
	for _, op := range []opcode.Opcode{
		regular(opcode.Op(255), 0, 0, 0, 0),
		regular(opcode.OpLoad, 9, 0, 0, register.ReturnsBase),
		regular(opcode.OpStore, 9, 0, 0, register.ReturnsBase),
		regular(opcode.OpLib, 0, 0, 0, 0),
		regular(opcode.OpDiv, 0, 1, 0, register.ReturnsBase),
	} {
		mustPanic(t, func() { New(op).Execute() })
	}
}

func TestNestedCodeCalls(t *testing.T) {
	leaf := New(
		regular(opcode.OpAdd, 0, register.ParamsBase, 1, register.ReturnsBase),
		immediate(999, register.TemporaryBase),
	)
	middle := New(
		immediate(100, register.TemporaryBase),
		regular(opcode.OpAdd, 0, register.ParamsBase, 0, register.ExtraBase),
		regular(opcode.OpLib, 0, 0, 0, 0),
		regular(opcode.OpAdd, 0, register.ExtraBase, register.TemporaryBase, register.ReturnsBase),
	).AddLib(leaf)
	p := New(
		immediate(500, register.TemporaryBase),
		immediate(41, register.ExtraBase),
		regular(opcode.OpLib, 0, 0, 0, 0),
		regular(opcode.OpAdd, 0, register.ExtraBase, 0, register.ReturnsBase),
		immediate(7, register.ExtraBase),
		regular(opcode.OpLib, 0, 0, 0, 0),
		regular(opcode.OpAdd, 0, register.ExtraBase, register.TemporaryBase, register.ReturnsBase+1),
	).AddLib(middle)
	p.Execute()
	if p.GetReturn(0) != 142 || p.GetReturn(1) != 608 || len(p.stack) != 1 || p.stack.GetTemporary(0) != 500 {
		t.Fatalf("nested/repeated call results = %d/%d, depth = %d", p.GetReturn(0), p.GetReturn(1), len(p.stack))
	}
	if middle.programCount != 0 || leaf.programCount != 0 || middle.stack.GetTemporary(0) != 0 || leaf.stack.GetTemporary(0) != 0 {
		t.Fatal("registered code source runtime was mutated")
	}
}

func TestCodeCallSharesMemory(t *testing.T) {
	lib := New(
		regular(opcode.OpStore, 1, register.ParamsBase, 0, register.ParamsBase+1),
		regular(opcode.OpAdd, 0, register.ParamsBase, 0, register.ReturnsBase),
	)
	p := New(
		regular(opcode.OpAlloc, 0, 1, 0, register.ExtraBase),
		immediate(65, register.ExtraBase+1),
		regular(opcode.OpLib, 0, 0, 0, 0),
		regular(opcode.OpLoad, 1, register.ExtraBase, 0, register.ReturnsBase),
	).AddLib(lib)
	p.Execute()
	if p.GetReturn(0) != 65 || p.memory.Get(0, 0) != 65 || lib.memory.Size(0) != 0 {
		t.Fatal("memory handle did not retain meaning across the frame call")
	}
}

func TestFailedCallRestoresCaller(t *testing.T) {
	lib := New(immediate(999, register.ReturnsBase), regular(opcode.Op(255), 0, 0, 0, 0))
	p := New(regular(opcode.OpLib, 0, 0, 0, 0), immediate(42, register.ReturnsBase)).AddLib(lib)
	p.stack.SetExtra(0, 7)
	mustPanic(t, p.Execute)
	if len(p.stack) != 1 || p.stack.GetExtra(0) != 7 || p.programCount != 1 || len(p.libs) != 1 {
		t.Fatal("failed child leaked a frame, partial returns or execution state")
	}
	p.Execute()
	if p.GetReturn(0) != 42 {
		t.Fatal("caller code was not restored")
	}
}

func TestArguments(t *testing.T) {
	p := New().AddArgs(42).AddStringArgs(String([]rune("첫째")), String([]rune("둘째"))).AddArgs(99)
	if p.stack.GetParams(0) != 42 || p.stack.GetParams(3) != 99 || p.paramCount != 4 {
		t.Fatal("argument indices did not advance")
	}
	for i, want := range []string{"첫째", "둘째"} {
		if got := string(ReadString(&p.memory, p.stack.GetParams(uint16(i+1)), 0)); got != want {
			t.Fatalf("argument %d = %q", i+1, got)
		}
	}
	p.paramCount = register.BankSize
	mustPanic(t, func() { p.AddArgs(1) })
	mustPanic(t, func() { p.AddStringArgs(String([]rune("overflow"))) })
	mustPanic(t, func() { p.GetReturn(register.BankSize) })
	mustPanic(t, func() { p.AddLib(nil) })
}
