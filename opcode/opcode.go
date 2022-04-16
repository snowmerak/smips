package opcode

// OpCode
type OpCode uint64

// Function is the function field.
func (op OpCode) Function() uint8 {
	return uint8(op & 0b11111111)
}

//Type
type Type uint8

// OpCodeTypes
const (
	TypeR = Type(0b00000000)
	TypeI = Type(0b10000010)
	TypeJ = Type(0b11110000)
)

// Type returns the type of the opcode.
func (op OpCode) Type() Type {
	t := Type(op & 0b11111111)
	if t >= TypeJ {
		return TypeJ
	}
	if t < TypeJ && t >= TypeI {
		return TypeI
	}
	return TypeR
}

// NewR creates a new R-type opcode.
func NewR(rsA, rsB, rd uint8, data uint32, funt uint8) R {
	return R(uint64(rsA)<<56 | uint64(rsB)<<48 | uint64(rd)<<40 | uint64(data)<<8 | uint64(funt))
}

// NewI creates a new I-type opcode.
func NewI(rs, rd uint8, imm uint32, funt uint8) I {
	return I(uint64(rs)<<56 | uint64(rd)<<48 | uint64(imm)<<8 | uint64(funt))
}

// NewJ creates a new J-type opcode.
func NewJ(address uint64, funt uint8) J {
	return J(uint64(address)<<8 | uint64(funt))
}
