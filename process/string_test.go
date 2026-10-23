package process

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/snowmerak/smips/memory"
)

func TestStringRoundTrip(t *testing.T) {
	for _, text := range []string{"", "A", "안녕🌏\x00!", strings.Repeat("가", 256)} {
		t.Run(textLengthName(text), func(t *testing.T) {
			m := memory.New()
			addr := WriteString(&m, String([]rune(text)))
			if got := string(ReadString(&m, addr, 0)); got != text {
				t.Fatalf("round trip = %q, want %q", got, text)
			}
			var header [8]byte
			binary.BigEndian.PutUint64(header[:], uint64(len([]rune(text))))
			for i, b := range header {
				if m.Get(addr, uint64(i)) != b {
					t.Fatal("length is not encoded as big-endian rune count")
				}
			}
		})
	}
}

func textLengthName(text string) string {
	if len(text) > 40 {
		return "multi-byte-length"
	}
	return text
}

func TestReadStringAtOffset(t *testing.T) {
	m := memory.New()
	addr := m.Alloc(15)
	// One rune U+1F30F, following a three-byte prefix.
	for i, b := range []byte{0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0xf3, 0x0f} {
		m.Set(addr, uint64(3+i), b)
	}
	if got := string(ReadString(&m, addr, 3)); got != "🌏" {
		t.Fatalf("decoded = %q", got)
	}
}

func TestReadStringRejectsInvalidBounds(t *testing.T) {
	for _, tc := range []struct {
		size, offset    uint64
		malformedLength bool
	}{
		{0, 0, false}, {7, 0, false}, {8, 1, false},
		{8, ^uint64(0), false}, {8, 0, true},
	} {
		m := memory.New()
		addr := m.Alloc(tc.size)
		if tc.malformedLength {
			for i := uint64(0); i < 8; i++ {
				m.Set(addr, i, 0xff)
			}
		}
		mustPanic(t, func() { ReadString(&m, addr, tc.offset) })
	}
}

func FuzzStringRoundTrip(f *testing.F) {
	for _, text := range []string{"", "A", "안녕🌏\x00", string([]byte{0xff, 0})} {
		f.Add(text)
	}
	f.Fuzz(func(t *testing.T, text string) {
		if len(text) > 4096 {
			t.Skip()
		}
		m := memory.New()
		want := String([]rune(text))
		addr := WriteString(&m, want)
		if got := ReadString(&m, addr, 0); string(got) != string(want) {
			t.Fatalf("round trip = %q, want %q", string(got), string(want))
		}
	})
}
