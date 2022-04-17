package memory

// Block is a memory block.
type Block struct {
	Value []uint8
	Count uint64
}

// Memory is a memory. lol.
type Memory struct {
	table map[uint64]*Block
	tombs []uint64
}

// New returns a new memory.
func New() Memory {
	return Memory{
		table: make(map[uint64]*Block),
		tombs: make([]uint64, 0),
	}
}

// Get is the getter of the memory.
// Address is the key of the memory.
// Offset is the index of the block.
func (m *Memory) Get(addr, offset uint64) uint8 {
	block, ok := m.table[addr]
	if !ok {
		return 0
	}
	if offset >= uint64(len(block.Value)) {
		return 0
	}
	return block.Value[offset]
}

// Set is the setter of the memory.
// Address is the key of the memory.
// Offset is the index of the block.
// Value is the value of the block.
func (m *Memory) Set(addr, offset uint64, value uint8) {
	block, ok := m.table[addr]
	if !ok {
		panic("memory: set: block not found")
	}
	if offset >= uint64(len(block.Value)) {
		panic("memory: set: offset out of range")
	}
	block.Value[offset] = value
}

// Free free the memory block at addr.
func (m *Memory) Free(addr uint64) uint8 {
	_, ok := m.table[addr]
	if !ok {
		return 0
	}
	m.tombs = append(m.tombs, addr)
	delete(m.table, addr)
	return 1
}

// Alloc allocates a memory block and return the address.
func (m *Memory) Alloc(size uint64) uint64 {
	addr := uint64(0)
	if len(m.tombs) > 0 {
		addr = m.tombs[len(m.tombs)-1]
		m.tombs = m.tombs[:len(m.tombs)-1]
	} else {
		addr = uint64(len(m.table))
	}
	m.table[addr] = &Block{
		Value: make([]uint8, size),
		Count: 1,
	}
	return addr
}
