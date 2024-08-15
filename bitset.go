package bitset

import (
	"bytes"
	"fmt"
	"math"
	"math/bits"
)

type Bitset struct {
	size int
	bits []uint64
}

func NewBitset(numBits int) *Bitset {
	numWords := int(math.Ceil(float64(numBits) / 64.0))
	return &Bitset{
		size: numBits,
		bits: make([]uint64, numWords),
	}
}

// Size returns the number of bits of the bitset
func (bitset *Bitset) Size() int {
	return bitset.size
}

// Set sets the Nth bit. Errors if n < 0 or n >= bitset.size
func (bitset *Bitset) Set(n int) error {
	if n < 0 {
		return fmt.Errorf("set: n must be >= 0")
	}
	if n >= bitset.size {
		return fmt.Errorf("bit index %d out of range of bitset", n)
	}
	idx := len(bitset.bits) - 1 - n/64
	bitset.bits[idx] |= 1 << (n % 64)
	return nil
}

// SetAll sets multiple bits.
func (bitset *Bitset) SetAll(bits ...int) error {
	for _, bit := range bits {
		if err := bitset.Set(bit); err != nil {
			return err
		}
	}
	return nil
}

// Clear zeroes the Nth bit. Errors if n < 0 or n >= bitset.size
func (bitset *Bitset) Clear(n int) error {
	if n < 0 {
		return fmt.Errorf("clear: n must be >= 0")
	}
	if n >= bitset.size {
		return fmt.Errorf("bit index %d out of range of bitset", n)
	}
	idx := len(bitset.bits) - 1 - n/64
	bitset.bits[idx] &= ^(1 << (n % 64))
	return nil
}

// ClearAll clears multiple bits.
func (bitset *Bitset) ClearAll(bits ...int) error {
	for _, bit := range bits {
		if err := bitset.Clear(bit); err != nil {
			return err
		}
	}
	return nil
}

// Flip flips the Nth bit, i.e. 0 -> 1 or 1 -> 0. Errors if n < 0 or n >= bitset.size
func (bitset *Bitset) Flip(n int) error {
	if n < 0 {
		return fmt.Errorf("clear: n must be >= 0")
	}
	if n >= bitset.size {
		return fmt.Errorf("bit index %d out of range of bitset", n)
	}
	idx := len(bitset.bits) - 1 - n/64
	bitset.bits[idx] ^= 1 << (n % 64)
	return nil
}

// FlipAll flips multiple bits.
func (bitset *Bitset) FlipAll(bits ...int) error {
	for _, bit := range bits {
		if err := bitset.Flip(bit); err != nil {
			return err
		}
	}
	return nil
}

// Test checks if the Nth bit is set. Errors if n < 0 or n >= bitset.size
func (bitset *Bitset) Test(n int) (bool, error) {
	if n < 0 {
		return false, fmt.Errorf("test: n must be >= 0")
	}
	if n >= bitset.size {
		return false, fmt.Errorf("bit index %d out of range of bitset", n)
	}
	idx := len(bitset.bits) - 1 - n/64
	return bitset.bits[idx]&(1<<(n%64)) >= 1, nil
}

func (bitset *Bitset) Not() {
	for i := range bitset.bits {
		bitset.bits[i] ^= bitset.bits[i]
	}
}

// Count returns the number of set bits
func (bitset *Bitset) Count() int {
	sum := 0
	for _, word := range bitset.bits {
		sum += bits.OnesCount64(word)
	}
	return sum
}

func (bitset *Bitset) String() string {
	buffer := bytes.NewBufferString("")
	for _, word := range bitset.bits {
		buffer.WriteString(fmt.Sprintf("%b", word))
	}
	return buffer.String()
}
