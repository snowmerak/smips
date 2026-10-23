package register

import "testing"

func TestRegisterAddressSpace(t *testing.T) {
	s := NewStack()
	s.Push(New())
	defer func() { s.Pop().Return() }()
	for i := 0; i < 1<<16; i++ {
		s.Set(uint16(i), uint64(i)+100)
	}
	for i := 0; i < 1<<16; i++ {
		if got := s.Get(uint16(i)); got != uint64(i)+100 {
			t.Fatalf("register %d = %d", i, got)
		}
	}
	// Bank-local access must name the same slots as instruction addressing.
	for i := 0; i < single; i++ {
		index := uint16(i)
		if s.GetParams(index) != uint64(4+i)+100 ||
			s.GetTemporary(index) != uint64(4+single+i)+100 ||
			s.GetExtra(index) != uint64(4+single*2+i)+100 ||
			s.GetReturns(index) != uint64(4+single*3+i)+100 {
			t.Fatalf("bank mapping differs at index %d", i)
		}
	}
	if s.GetHi() != s.Get(uint16(4+single*2-1)) {
		t.Fatal("Hi does not alias the last temporary register")
	}
}

func TestNewResetsAllBanks(t *testing.T) {
	r := New()
	for i := 0; i < single; i++ {
		r.params[i], r.temporary[i], r.extra[i], r.returns[i] = 9, 9, 9, 9
	}
	r.zero, r.one, r.two, r.four = 9, 9, 9, 9
	r.Return()
	r = New()
	defer r.Return()
	if r.zero != 0 || r.one != 1 || r.two != 2 || r.four != 4 {
		t.Fatal("reserved register initial values were not restored")
	}
	for i := 0; i < single; i++ {
		if r.params[i] != 0 || r.temporary[i] != 0 || r.extra[i] != 0 || r.returns[i] != 0 {
			t.Fatalf("reused register contains stale data at %d", i)
		}
	}
}

func TestNestedFrames(t *testing.T) {
	s := NewStack()
	s.Push(New())
	defer func() { s.Pop().Return() }()
	s.SetTemporary(0, 100)
	s.SetExtra(0, 42)
	s.SetExtra(single-1, 77)
	s.Call(New())
	if s.GetParams(0) != 42 || s.GetParams(single-1) != 77 || s.GetTemporary(0) != 0 {
		t.Fatal("child input or register isolation is incorrect")
	}
	s.SetTemporary(0, 200)
	s.SetExtra(0, 7)
	s.Call(New())
	s.SetReturns(0, s.GetParams(0)+1)
	s.Plunder()
	if s.GetTemporary(0) != 200 || s.GetExtra(0) != 8 {
		t.Fatal("nested return corrupted the caller")
	}
	s.SetReturns(0, s.GetParams(0)+s.GetExtra(0))
	s.SetReturns(single-1, 99)
	s.Plunder()
	if len(s) != 1 || s.GetTemporary(0) != 100 || s.GetExtra(0) != 50 || s.GetExtra(single-1) != 99 {
		t.Fatal("return did not restore the root frame and export results")
	}
}

func TestPlunderPreservesRoot(t *testing.T) {
	s := NewStack()
	s.Push(New())
	defer func() { s.Pop().Return() }()
	defer func() {
		if recover() == nil {
			t.Error("root frame return must panic")
		}
		if len(s) != 1 {
			t.Error("failed return removed the root frame")
		}
	}()
	s.Plunder()
}
