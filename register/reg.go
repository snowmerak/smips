package register

import "sync"

// BankSize is the number of registers in each frame-local bank.
const BankSize = 16383

// Instruction addresses cover four reserved slots followed by the four banks.
const (
	ParamsBase    uint16 = 4
	TemporaryBase        = ParamsBase + BankSize
	ExtraBase            = TemporaryBase + BankSize
	ReturnsBase          = ExtraBase + BankSize
)

const single = BankSize

// Register is the register of the VM.
// Addresses 0, 1, 2, 3 initially contain uint64 values 0, 1, 2, 4.
// 16383 registers of params, 16383 registers of temporary,
// 16383 registers of extra, 16383 registers of returns.
// extra is the parameter of stack when call, returns is the return of stack when plunder.
// when last stack is plunder, the extra is the return of last stack.
type Register struct {
	zero      uint64
	one       uint64
	two       uint64
	four      uint64
	params    [single]uint64
	temporary [single]uint64
	extra     [single]uint64
	returns   [single]uint64
}

// New is the constructor of Register.
func New() *Register {
	reg := pool.Get().(*Register)
	*reg = Register{one: 1, two: 2, four: 4}
	return reg
}

// Return is the destructor of Register.
func (r *Register) Return() {
	pool.Put(r)
}

// pool is the pool of Register.
var pool = sync.Pool{
	New: func() any {
		return &Register{}
	},
}

// Stack is the stack of Register.
type Stack []*Register

// NewStack is the constructor of Stack.
func NewStack() Stack {
	return Stack{}
}

// Last is the last register of Stack.
func (s *Stack) Last() *Register {
	return (*s)[len(*s)-1]
}

// Push is the push of Stack.
func (s *Stack) Push(r *Register) {
	*s = append(*s, r)
}

// Call is append the register to stack.
// the pushed register will get the params of the last register via extra.
func (s *Stack) Call(reg *Register) {
	last := s.Last()
	for i := 0; i < single; i++ {
		reg.params[i] = last.extra[i]
	}
	s.Push(reg)
}

// Pop is the pop of Stack.
func (s *Stack) Pop() *Register {
	r := (*s)[len(*s)-1]
	(*s)[len(*s)-1] = nil
	*s = (*s)[:len(*s)-1]
	return r
}

// Plunder remove the last register of Stack.
// the removed register will export the returns to the extra of the last register.
func (s *Stack) Plunder() {
	if len(*s) < 2 {
		panic("register: cannot return from root frame")
	}
	pop := s.Pop()
	last := s.Last()
	for i := 0; i < single; i++ {
		last.extra[i] = pop.returns[i]
	}
	pop.Return()
}

// Set is the set of Register.
func (r *Stack) Set(index uint16, value uint64) {
	*r.Last().slot(index) = value
}

// Get is the get of Register.
func (r *Stack) Get(index uint16) uint64 {
	return *r.Last().slot(index)
}

func (r *Register) slot(index uint16) *uint64 {
	switch index {
	case 0:
		return &r.zero
	case 1:
		return &r.one
	case 2:
		return &r.two
	case 3:
		return &r.four
	default:
		if index < TemporaryBase {
			return &r.params[index-ParamsBase]
		} else if index < ExtraBase {
			return &r.temporary[index-TemporaryBase]
		} else if index < ReturnsBase {
			return &r.extra[index-ExtraBase]
		} else {
			return &r.returns[index-ReturnsBase]
		}
	}
}

// SetParams is the set of params of Register.
func (s *Stack) SetParams(index uint16, value uint64) {
	r := s.Last()
	if index < single {
		r.params[index] = value
	}
}

// GetParams is the get of params of Register.
func (s *Stack) GetParams(index uint16) uint64 {
	r := s.Last()
	if index < single {
		return r.params[index]
	}
	return 0
}

// SetTemporary is the set of params of Register.
func (s *Stack) SetTemporary(index uint16, value uint64) {
	r := s.Last()
	if index < single {
		r.temporary[index] = value
	}
}

// GetTemporary is the get of params of Register.
func (s *Stack) GetTemporary(index uint16) uint64 {
	r := s.Last()
	if index < single {
		return r.temporary[index]
	}
	return 0
}

// SetHi is the set of params of Register.
func (s *Stack) SetHi(value uint64) {
	r := s.Last()
	r.temporary[single-1] = value
}

// GetHi is the get of params of Register.
func (s *Stack) GetHi() uint64 {
	r := s.Last()
	return r.temporary[single-1]
}

// SetExtra is the set of params of Register.
func (s *Stack) SetExtra(index uint16, value uint64) {
	r := s.Last()
	if index < single {
		r.extra[index] = value
	}
}

// GetExtra is the get of params of Register.
func (s *Stack) GetExtra(index uint16) uint64 {
	r := s.Last()
	if index < single {
		return r.extra[index]
	}
	return 0
}

// SetReturns is the set of params of Register.
func (s *Stack) SetReturns(index uint16, value uint64) {
	r := s.Last()
	if index < single {
		r.returns[index] = value
	}
}

// GetReturns is the get of params of Register.
func (s *Stack) GetReturns(index uint16) uint64 {
	r := s.Last()
	if index < single {
		return r.returns[index]
	}
	return 0
}
