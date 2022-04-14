package opcode

// I
// 8bit rs, 8bit rd, 40bit imm, 8bit funct
type I uint64

// ResisterSouce is the source register.
func (i I) RegisterSource() uint8 {
	return uint8(i >> 56 & 0b11111111)
}

// RegisterDestination is the destination register.
func (i I) RegisterDestination() uint8 {
	return uint8(i >> 48 & 0b11111111)
}

// Immediate is the immediate value.
func (i I) Immediate() uint64 {
	return uint64(i >> 8 & 0b1111111111111111111111111111111111111111)
}

// Function is the function field.
func (i I) Function() uint8 {
	return uint8(i & 0b11111111)
}
