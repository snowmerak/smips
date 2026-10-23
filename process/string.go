package process

import (
	"github.com/snowmerak/smips/memory"
)

// String is a sequence of runes.
// in memory, first 8 bytes is length of string.
// and 4 bytes is each rune.
type String []rune

// ReadString reads string from memory.
func ReadString(m *memory.Memory, pos, off uint64) String {
	size := m.Size(pos)
	if off > size || size-off < 8 {
		panic("process: string header out of range")
	}
	length := uint64(0)
	for i := uint64(0); i < 8; i++ {
		b := uint64(m.Get(pos, off+i))
		length = (length << 8) | b
	}
	// Check before allocation so an invalid header cannot request an unbounded slice.
	if length > (size-off-8)/4 {
		panic("process: string contents out of range")
	}
	runes := make([]rune, length)
	for i := uint64(0); i < length; i++ {
		v := int32(0)
		for j := uint64(0); j < 4; j++ {
			b := int32(m.Get(pos, off+8+i*4+j))
			v = (v << 8) | b
		}
		runes[i] = v
	}
	return runes
}

// WriteString writes string to memory.
func WriteString(m *memory.Memory, s String) uint64 {
	length := uint64(len(s))
	addr := m.Alloc(8 + length*4)
	for i := 0; i < 8; i++ {
		m.Set(addr, uint64(7-i), uint8(length>>(8*uint(i))))
	}
	for i := uint64(0); i < length; i++ {
		v := uint32(s[i])
		for j := 0; j < 4; j++ {
			m.Set(addr, 8+i*4+uint64(3-j), uint8(v>>(8*uint(j))))
		}
	}
	return addr
}
