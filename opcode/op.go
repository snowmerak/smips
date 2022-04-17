package opcode

// Op is a operation of instruction.
type Op uint8

// Op List.
const (
	// OpAdd is add operation.
	// rd = ra + rb
	OpAdd Op = iota
	// OpSub is sub operation.
	// rd = ra - rb
	OpSub
	// OpMul is mul operation.
	// rd = ra * rb
	OpMul
	// OpDiv is div operation.
	// rd = ra / rb
	OpDiv
	// OpMod is mod operation.
	// rd = ra % rb
	OpMod
	// OpAnd is and operation.
	// rd = ra & rb
	OpAnd
	// OpOr is or operation.
	// rd = ra | rb
	OpOr
	// OpXor is xor operation.
	// rd = ra ^ rb
	OpXor
	// OpNot is not operation.
	// rd = ~ra
	OpNot
	// OpShl is shl operation.
	// rd = ra << rb
	OpShl
	// OpShr is shr operation.
	// rd = ra >> rb
	OpShr
	// OpEq is eq operation.
	// rd = ra == rb
	OpEq
	// OpNeq is ne operation.
	// rd = ra != rb
	OpNeq
	// OpFAdd is lt operation.
	// rd = ra + rb in float
	OpFAdd
	// OpFSub is lt operation.
	// rd = ra - rb in float
	OpFSub
	// OpFMul is lt operation.
	// rd = ra * rb in float
	OpFMul
	// OpFDiv is lt operation.
	// rd = ra / rb in float
	OpFDiv
	// OpFMod is lt operation.
	// rd = ra % rb in float
	OpFMod
	// OpSet is lt operation.
	// rd = wide data
	OpSet
	// OpJmp is jmp operation.
	// pc = ra
	OpJmp
	// OpAlloc is alloc operation.
	// ra is size of memory block
	// rd is address of memory block
	OpAlloc
	// OpLoad is load operation.
	// rd = ra[rb] with data(size).
	// data is 8 bit data.
	OpLoad
	// OpStore is store operation.
	// ra[rb] = rd with data(size).
	// data is 3 bit size for this
	OpStore
	// OpLib is library call operation.
	// ra is address of library.
	// can pass prameter to library via extra registers.
	// can receive return value from library via extra registers.
	OpLib

	// --------------------
	// System Calls.
	// data is extra data.
	// ra, rb is source of system call.
	// rd is destination of system call.

	// OpenFile
	// Ra is String address of file name on memory.
	// Rb is Mode.
	// Rd will get the file descriptor.
	OpenFile
	// CloseFile
	// Ra is the file descriptor.
	// Rd will get the closed status.
	CloseFile
	// ReadFile
	// Ra is the file descriptor.
	// Data is size.
	// Rd will get the read bytes.
	ReadFile
	// ReadFileAt
	// Ra is the file descriptor.
	// Rb is the offset.
	// Data is size.
	// Rd will get the read bytes.
	ReadFileAt
	// WriteFile
	// Ra is the file descriptor.
	// Data is size.
	// Rd is the value.
	// Hi will get the written bytes.
	WriteFile
	// WriteFileAt
	// Ra is the file descriptor.
	// Rb is the offset.
	// Data is size.
	// Rd is the value.
	// Hi will get the written bytes.
	WriteFileAt
)
