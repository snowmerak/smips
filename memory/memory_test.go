package memory

import "testing"

func TestAllocationReuse(t *testing.T) {
	m := New()
	a, b := m.Alloc(2), m.Alloc(3)
	m.Set(a, 1, 65)
	m.Set(b, 2, 66)
	if m.Size(a) != 2 || m.Get(a, 1) != 65 || m.Get(b, 2) != 66 {
		t.Fatal("blocks are not independent")
	}
	if m.Free(a) != 1 || m.Free(a) != 0 || m.Get(a, 1) != 0 || m.Size(a) != 0 {
		t.Fatal("freed block is still accessible or freed twice")
	}
	c := m.Alloc(1)
	if c != a || m.Size(c) != 1 || m.Get(c, 0) != 0 || m.Get(b, 2) != 66 {
		t.Fatal("reallocation did not clear the block or corrupted another block")
	}
	d := m.Alloc(0)
	if d == b || d == c || m.Size(d) != 0 || m.Get(d, 0) != 0 {
		t.Fatal("zero-sized allocation replaced an existing block")
	}
}

func TestMemoryBounds(t *testing.T) {
	m := New()
	a := m.Alloc(1)
	if m.Get(a, 1) != 0 || m.Get(99, 0) != 0 {
		t.Fatal("invalid reads must return zero")
	}
	for _, addr := range []uint64{a, 99} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("invalid write must panic")
				}
			}()
			m.Set(addr, 1, 1)
		}()
	}
}
