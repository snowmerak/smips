package opcode

// Opcode is a single instruction in the VM.
// 8 bit data, 16 bit ra, 16 bit rb, 16 bit rd, 8 bit op.
// 40 bit wide-data, 16 bit rd, 8 bit op.
type Opcode uint64

func NewRegularOpCode(data uint8, ra, rb, rd uint16, op uint8) Opcode {
	return Opcode((uint64(data) << 56) | (uint64(ra) << 40) | (uint64(rb) << 24) | (uint64(rd) << 8) | uint64(op))
}

func NewWideOpCode(data uint64, rd uint16, op uint8) Opcode {
	if data >= 1<<40 {
		panic("opcode: immediate exceeds 40 bits")
	}
	return Opcode((uint64(data) << 24) | (uint64(rd) << 8) | uint64(op))
}

// Data is the 8 bit data field of the opcode.
func (o Opcode) Data() uint8 {
	return uint8(o >> 56)
}

// Ra is the 16 bit ra field of the opcode.
func (o Opcode) Ra() uint16 {
	return uint16(o >> 40)
}

// Rb is the 16 bit rb field of the opcode.
func (o Opcode) Rb() uint16 {
	return uint16(o >> 24)
}

// WideData is the 40 bit wide-data field of the opcode.
func (o Opcode) WideData() uint64 {
	return uint64(o >> 24)
}

// Rd is the 16 bit rd field of the opcode.
func (o Opcode) Rd() uint16 {
	return uint16(o >> 8)
}

// Op is the 8 bit op field of the opcode.
func (o Opcode) Op() uint8 {
	return uint8(o)
}

// Size interprets Data as a byte count for load/store instructions.
// A uint64 holds at most eight bytes; zero is a no-op.
func (o Opcode) Size() uint8 {
	size := o.Data()
	if size > 8 {
		panic("opcode: byte count exceeds 8")
	}
	return size
}
