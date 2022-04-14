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
	TypeR = Type(0b11111101)
	TypeI = Type(0b11111110)
	TypeJ = Type(0b11111111)
)

// Type returns the type of the opcode.
func (op OpCode) Type() Type {
	t := Type(op & 0b11111111)
	if t == TypeJ {
		return TypeJ
	}
	if t == TypeI {
		return TypeI
	}
	return TypeR
}