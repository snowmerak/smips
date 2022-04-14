package register

import "sync"

// Register is a register. lol.
type Register struct {
	zero uint64      // 0
	at   uint64      // 1
	v    [32]uint64  // 2-33
	a    [32]uint64  // 34-65
	t    [100]uint64 // 66-165
	s    [70]uint64  // 166-235
	k    [16]uint64  // 236-251
	gp   uint64      // 252
	sp   uint64      // 253
	fp   uint64      // 254
	ra   uint64      // 255
}

var pool = sync.Pool{
	New: func() interface{} {
		return &Register{}
	},
}

// New returns a new Register.
func New() *Register {
	return pool.Get().(*Register)
}

// Return returns register to the pool.
func (r *Register) Return() {
	for i := 0; i < 256; i++ {
		r.Set(byte(i), 0)
	}
	pool.Put(r)
}

// Set sets the value of the register at the given index.
func (r *Register) Set(index byte, value uint64) {
	if index == 0 {
		r.zero = value
		return
	}
	if index == 1 {
		r.at = value
		return
	}
	if index <= 33 {
		r.v[index-2] = value
		return
	}
	if index <= 65 {
		r.a[index-34] = value
		return
	}
	if index <= 165 {
		r.t[index-66] = value
		return
	}
	if index <= 235 {
		r.s[index-166] = value
		return
	}
	if index <= 251 {
		r.k[index-236] = value
		return
	}
	if index == 252 {
		r.gp = value
		return
	}
	if index == 253 {
		r.sp = value
		return
	}
	if index == 254 {
		r.fp = value
		return
	}
	r.ra = value
}

// Get returns the value of the register at the given index.
func (r *Register) Get(index byte) uint64 {
	if index == 0 {
		return r.zero
	}
	if index == 1 {
		return r.at
	}
	if index <= 33 {
		return r.v[index-2]
	}
	if index <= 65 {
		return r.a[index-34]
	}
	if index <= 165 {
		return r.t[index-66]
	}
	if index <= 235 {
		return r.s[index-166]
	}
	if index <= 251 {
		return r.k[index-236]
	}
	if index == 252 {
		return r.gp
	}
	if index == 253 {
		return r.sp
	}
	if index == 254 {
		return r.fp
	}
	return r.ra
}

// GetV returns the value of the register at the given index. index <= 32
// value of function returns.
func (r *Register) GetV(index byte) uint64 {
	if index > 32 {
		return 0
	}
	return r.v[index]
}

// SetV sets the value of the register at the given index. index <= 32
// value of function returns.
func (r *Register) SetV(index byte, value uint64) {
	if index > 32 {
		return
	}
	r.v[index] = value
}

// GetA returns the value of the register at the given index. index <= 32
// value of function paramters.
func (r *Register) GetA(index byte) uint64 {
	if index > 32 {
		return 0
	}
	return r.a[index]
}

// SetA sets the value of the register at the given index. index <= 32
// value of function paramters.
func (r *Register) SetA(index byte, value uint64) {
	if index > 32 {
		return
	}
	r.a[index] = value
}

// GetT returns the value of the register at the given index. index <= 100
// temporary values.
func (r *Register) GetT(index byte) uint64 {
	if index > 100 {
		return 0
	}
	return r.t[index]
}

// SetT sets the value of the register at the given index. index <= 100
// temporary values.
func (r *Register) SetT(index byte, value uint64) {
	if index > 100 {
		return
	}
	r.t[index] = value
}

// GetHi returns the value of the register.
// high 64 bits of 128 bit value.
func (r *Register) GetHi() uint64 {
	return r.t[len(r.t)-1]
}

// SetHi sets the value of the register.
// high 64 bits of 128 bit value.
func (r *Register) SetHi(value uint64) {
	r.t[len(r.t)-1] = value
}

// GetS returns the value of the register at the given index. index <= 70
// persistance values.
func (r *Register) GetS(index byte) uint64 {
	if index > 70 {
		return 0
	}
	return r.s[index]
}

// SetS sets the value of the register at the given index. index <= 70
// persistance values.
func (r *Register) SetS(index byte, value uint64) {
	if index > 70 {
		return
	}
	r.s[index] = value
}

// GetK returns the value of the register at the given index. index <= 16
// values for os kernels.
func (r *Register) GetK(index byte) uint64 {
	if index > 16 {
		return 0
	}
	return r.k[index]
}

// SetK sets the value of the register at the given index. index <= 16
// values for os kernels.
func (r *Register) SetK(index byte, value uint64) {
	if index > 16 {
		return
	}
	r.k[index] = value
}

// GetGP returns the value of the register.
// global poiner.
func (r *Register) GetGP() uint64 {
	return r.gp
}

// SetGP sets the value of the register.
// global poiner.
func (r *Register) SetGP(value uint64) {
	r.gp = value
}

// GetSP returns the value of the register.
// stack pointer.
func (r *Register) GetSP() uint64 {
	return r.sp
}

// SetSP sets the value of the register.
// stack pointer.
func (r *Register) SetSP(value uint64) {
	r.sp = value
}

// GetFP returns the value of the register.
// function pointer.
func (r *Register) GetFP() uint64 {
	return r.fp
}

// SetFP sets the value of the register.
// function pointer.
func (r *Register) SetFP(value uint64) {
	r.fp = value
}

// GetRA returns the value of the register.
// returning address.
func (r *Register) GetRA() uint64 {
	return r.ra
}

// SetRA sets the value of the register.
// returning address.
func (r *Register) SetRA(value uint64) {
	r.ra = value
}
