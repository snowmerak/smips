package opcode

// R
// 8bit rs-a, 8bit rs-b, 8bit rd, 32(8x4) data, 8bit funct,
type R uint64

// RegisterSourceA is the first source register.
func (r R) RegisterSourceA() uint8 {
	return uint8(r >> 56 & 0b11111111)
}

// RegisterSourceB is the second source register.
func (r R) RegisterSourceB() uint8 {
	return uint8(r >> 48 & 0b11111111)
}

// RegisterDestination is the destination register.
func (r R) RegisterDestination() uint8 {
	return uint8(r >> 40 & 0b11111111)
}

// Data is the data field.
func (r R) Data() uint32 {
	return uint32(r >> 8 & 0b11111111111111111111111111111111)
}

// Function is the function field.
func (r R) Function() uint8 {
	return uint8(r & 0b11111111)
}
