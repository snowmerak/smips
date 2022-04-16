package opcode

// Funcs are opcode functions.
const (
	RAdd = uint8(TypeR + iota)
	RSub
	RMul
	RDiv
	RRem
	RAnd
	ROr
	RXor
	RNor
	RShiftLeft
	RShiftRightLogical
	RShiftRightArithmetic
	RBranchEqual
	RBranchNotEqual
	RBranchLessThan
	RBranchLessThanOrEqual
	RBranchGreaterThan
	RBranchGreaterThanOrEqual
	RFloatAdd
	RFloatSub
	RFloatMul
	RFloatDiv
)
const (
	ISet = uint8(TypeI + iota)
	IAdd
	ISub
	IMul
	IDiv
	IRem
	ILoadData
	IStoreData
)
const (
	Jump = uint8(TypeJ + iota)
	JumpTo
	JumpAndAddRegister
	JumpAndPassParameter
	JumpAndRemoveRegister
	JumpAndReturnValue
	JPushRegister
	JBringParameter
	JPopRegister
	JBringReturns
)
