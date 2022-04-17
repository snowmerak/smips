package register

import "sync"

// single is the number of each registers.
const single = 16383

// Register is the register of the VM.
// 8 bit zero, 8 bit one, 8 bit two, 8 bit four,
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
	reg.zero = 0
	reg.one = 1
	reg.two = 2
	reg.four = 4
	for i := 0; i < single; i++ {
		reg.params[i] = 0
	}
	for i := 0; i < single*2; i++ {
		reg.temporary[i] = 0
	}
	for i := 0; i < single; i++ {
		reg.returns[i] = 0
	}
	return reg
}

// Return is the destructor of Register.
func (r *Register) Return() {
	pool.Put(r)
}

// pool is the pool of Register.
var pool = sync.Pool{
	New: func() interface{} {
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
	*s = (*s)[:len(*s)-1]
	return r
}

// Plunder remove the last register of Stack.
// the removed register will export the returns to the extra of the last register.
func (s *Stack) Plunder() {
	pop := s.Pop()
	last := s.Last()
	for i := 0; i < single; i++ {
		last.extra[i] = pop.returns[i]
	}
	pop.Return()
}
