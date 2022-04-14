package process

import (
	"github.com/snowmerak/smips/memory"
	"github.com/snowmerak/smips/stack"
)

type Process struct {
	memory *memory.Memory
	stack  *stack.Stack
}
