package opcode

// J
// 58bit address, 8bit funct
type J uint64

func (j J) Address() uint64 {
	return uint64(j >> 8 & 0b1111111111111111111111111111111111111111111111111111111111111111)
}

func (j J) Function() uint8 {
	return uint8(j & 0b11111111)
}
