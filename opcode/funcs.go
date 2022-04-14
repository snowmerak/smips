package opcode

// Funcs are opcode functions.
const (
	RAdd = uint8(TypeR + iota)
)
const (
	IAdd = uint8(TypeI + iota)
)
const (
	Jump = uint8(TypeJ + iota)
	JumpAndAddRegister
	JumpAndRemoveRegister
)
