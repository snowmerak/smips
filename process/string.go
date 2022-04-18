package process

import (
	"github.com/snowmerak/smips/memory"
)

// String is alias of []rune
// in memory, first 8 bytes is length of string.
// and 4 bytes is each rune.
type String []rune

// ReadString reads string from memory.
func ReadString(m *memory.Memory, pos, off uint64) String {
	len := uint64(0)
	for i := uint64(0); i < 8; i++ {
		b := uint64(m.Get(pos, off+i))
		if b == 0 {
			return []rune{}
		}
		len = (len << 8) | b
	}
	runes := make([]rune, len)
	for i := uint64(0); i < len; i++ {
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
	len := uint64(len(s))
	addr := m.Alloc(8 + len*4)
	offset := uint64(0)
	for i := 0; i < 8; i++ {
		m.Set(addr, 7-offset, uint8(len>>(8*uint64(i))))
	}
	offset += 8
	for i := uint64(0); i < len; i++ {
		v := uint32(s[i])
		for j := 0; j < 4; j++ {
			m.Set(addr, 7-offset, uint8(v>>(8*uint64(j))))
		}
		offset += 4
	}
	return addr
}
