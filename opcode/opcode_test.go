package opcode

import "testing"

func TestRegularEncoding(t *testing.T) {
	o := NewRegularOpCode(0x12, 0x3456, 0x789a, 0xbcde, 0xf0)
	if uint64(o) != 0x123456789abcdef0 {
		t.Fatalf("instruction = %016x", o)
	}
	// Decode an independent bit vector as well as checking the encoder.
	o = Opcode(0x123456789abcdef0)
	if o.Data() != 0x12 || o.Ra() != 0x3456 || o.Rb() != 0x789a || o.Rd() != 0xbcde || o.Op() != 0xf0 {
		t.Fatalf("decoded fields = %x %x %x %x %x", o.Data(), o.Ra(), o.Rb(), o.Rd(), o.Op())
	}
}

func TestWideEncoding(t *testing.T) {
	o := NewWideOpCode(0x123456789a, 0xbcde, 0xf0)
	if uint64(o) != 0x123456789abcdef0 {
		t.Fatalf("instruction = %016x", o)
	}
	o = Opcode(0x123456789abcdef0)
	if o.WideData() != 0x123456789a || o.Rd() != 0xbcde || o.Op() != 0xf0 {
		t.Fatalf("decoded fields = %x %x %x", o.WideData(), o.Rd(), o.Op())
	}
}

func TestWideEncodingRejectsOverflow(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("40-bit immediate overflow was silently truncated")
		}
	}()
	NewWideOpCode(1<<40, 0, 0)
}

func FuzzRegularEncoding(f *testing.F) {
	f.Add(uint8(0x12), uint16(0x3456), uint16(0x789a), uint16(0xbcde), uint8(0xf0))
	f.Add(uint8(255), uint16(65535), uint16(65535), uint16(65535), uint8(255))
	f.Fuzz(func(t *testing.T, data uint8, ra, rb, rd uint16, op uint8) {
		o := NewRegularOpCode(data, ra, rb, rd, op)
		if o.Data() != data || o.Ra() != ra || o.Rb() != rb || o.Rd() != rd || o.Op() != op {
			t.Fatalf("fields do not round trip: %016x", o)
		}
	})
}

func FuzzWideEncoding(f *testing.F) {
	f.Add(uint64(0x123456789a), uint16(0xbcde), uint8(0xf0))
	f.Add(uint64(1<<40-1), uint16(65535), uint8(255))
	f.Fuzz(func(t *testing.T, data uint64, rd uint16, op uint8) {
		data &= 1<<40 - 1
		o := NewWideOpCode(data, rd, op)
		if o.WideData() != data || o.Rd() != rd || o.Op() != op {
			t.Fatalf("fields do not round trip: %016x", o)
		}
	})
}
