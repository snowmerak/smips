package opcode

// Opcode is a single instruction in the VM.
// 8 bit data, 16 bit ra, 16 bit rb, 16 bit rd, 8 bit op.
// 40 bit wide-data, 16 bit rd, 8 bit op.
type Opcode uint64

func NewRegularOpCode(data uint8, ra, rb, rd uint16, op uint8) Opcode {
	return Opcode((uint64(data) << 56) | (uint64(ra) << 40) | (uint64(rb) << 24) | (uint64(rd) << 8) | uint64(op))
}

func NewWideOpCode(data uint64, rd uint16, op uint8) Opcode {
	return Opcode((uint64(data) << 24) | (uint64(rd) << 8) | uint64(op))
}

// Data is the 8 bit data field of the opcode.
func (o Opcode) Data() uint8 {
	return uint8(o >> 56)
}

// Ra is the 16 bit ra field of the opcode.
func (o Opcode) Ra() uint16 {
	return uint16(o >> 48)
}

// Rb is the 16 bit rb field of the opcode.
func (o Opcode) Rb() uint16 {
	return uint16(o >> 32)
}

// WideData is the 40 bit wide-data field of the opcode.
func (o Opcode) WideData() uint64 {
	return uint64(o >> 14)
}

// Rd is the 16 bit rd field of the opcode.
func (o Opcode) Rd() uint16 {
	return uint16(o >> 16)
}

// Op is the 8 bit op field of the opcode.
func (o Opcode) Op() uint8 {
	return uint8(o >> 8)
}
