package memory

import (
	"runtime"
	"sort"
	"sync/atomic"
)

// pair is a pair of start and end. start <= index < end.
type pair struct {
	start uint64
	end   uint64
}

// Memory is a memory. lol.
type Memory struct {
	space []uint64
	marks []pair
	lock  int64
}

// New returns a new Memory.
func New(size int) *Memory {
	return &Memory{
		space: make([]uint64, size*8),
		marks: make([]pair, 0),
	}
}

// Lock locks the memory.
func (m *Memory) Lock() {
	for atomic.CompareAndSwapInt64(&m.lock, 0, 1) {
		runtime.Gosched()
	}
}

// Unlock unlocks the memory.
func (m *Memory) Unlock() {
	atomic.StoreInt64(&m.lock, 0)
}

// GetFreeSpace returns the first free space.
func (m *Memory) GetFreeSpace(size uint64) uint64 {
	prev := uint64(0)
	for _, v := range m.marks {
		if v.start == 0 {
			continue
		}
		if v.start-prev >= size {
			return prev
		}
		prev = v.end
	}

	alloced := uint64(len(m.space)) - prev
	m.space = append(m.space, make([]uint64, size-alloced+1)...)
	return prev
}

// Alloc allocates a new space.
func (m *Memory) Alloc(size uint64) (uint64, bool) {
	m.Lock()
	defer m.Unlock()
	if size == 0 {
		return 0, false
	}

	address := m.GetFreeSpace(size)
	m.marks = append(m.marks, pair{address, address + size})
	return address, true
}

// Free frees a space.
func (m *Memory) Free(address uint64) bool {
	m.Lock()
	defer m.Unlock()
	idx := sort.Search(len(m.marks), func(i int) bool {
		return m.marks[i].start >= address
	})
	if idx == len(m.marks) || m.marks[idx].start != address {
		return false
	}
	s := m.marks[idx].start
	e := m.marks[idx].end
	for s < e {
		m.space[s] = 0
		s++
	}
	m.marks = append(m.marks[:idx], m.marks[idx+1:]...)
	return true
}

// Store stores a value.
func (m *Memory) Store(address uint64, value uint64) {
	m.space[address] = value
}

// Load loads a value.
func (m *Memory) Load(address uint64) uint64 {
	return m.space[address]
}

// End returns the end of the structrue.
func (m *Memory) End(address uint64) (uint64, bool) {
	idx := sort.Search(len(m.marks), func(i int) bool {
		return m.marks[i].start >= address
	})
	if idx == len(m.marks) || m.marks[idx].start != address {
		return 0, false
	}
	return m.marks[idx].end, true
}
