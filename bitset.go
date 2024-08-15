package bitset

import (
	"bytes"
	"fmt"
	"math"
	"math/bits"
)

type BitSet struct {
	size     int // the number of bits the bitset holds
	bitArray []uint64
}

// NewBitSet initializes and returns a BitSet with the given number of bits
func NewBitSet(numBits int) *BitSet {
	numWords := int(math.Ceil(float64(numBits) / 64.0))
	return &BitSet{
		size:     numBits,
		bitArray: make([]uint64, numWords),
	}
}

// Size returns the number of bits of the bitset
func (bitset *BitSet) Size() int {
	return bitset.size
}

// Set sets the Nth bit. Errors if n < 0 or n >= bitset.size
func (bitset *BitSet) Set(n int) error {
	if err := bitset.checkValidBit(n); err != nil {
		return err
	}
	bit := len(bitset.bitArray) - 1 - n/64
	bitset.bitArray[bit] |= 1 << (n % 64)
	return nil
}

// SetBits sets multiple bits.
func (bitset *BitSet) SetBits(bits []int) error {
	for _, bit := range bits {
		if err := bitset.Set(bit); err != nil {
			return err
		}
	}
	return nil
}

// Clear zeroes the Nth bit. Errors if n < 0 or n >= bitset.size
func (bitset *BitSet) Clear(n int) error {
	if err := bitset.checkValidBit(n); err != nil {
		return err
	}
	bit := len(bitset.bitArray) - 1 - n/64
	bitset.bitArray[bit] &= ^(1 << (n % 64))
	return nil
}

// ClearBits clears the bits at the given positions.
func (bitset *BitSet) ClearBits(bits []int) error {
	for _, bit := range bits {
		if err := bitset.Clear(bit); err != nil {
			return err
		}
	}
	return nil
}

// ClearAll clears all bits.
func (bitset *BitSet) ClearAll() {
	bitset.bitArray = make([]uint64, int(math.Ceil(float64(bitset.size)/64.0)))
}

// Flip flips the Nth bit, i.e. 0 -> 1 or 1 -> 0. Errors if n < 0 or n >= bitset.size
func (bitset *BitSet) Flip(n int) error {
	if err := bitset.checkValidBit(n); err != nil {
		return err
	}
	bit := len(bitset.bitArray) - 1 - n/64
	bitset.bitArray[bit] ^= 1 << (n % 64)
	return nil
}

// FlipBits flips multiple bits.
func (bitset *BitSet) FlipBits(bits []int) error {
	for _, bit := range bits {
		if err := bitset.Flip(bit); err != nil {
			return err
		}
	}
	return nil
}

// Test checks if the Nth bit is set. Errors if n < 0 or n >= bitset.size
func (bitset *BitSet) Test(n int) (bool, error) {
	if err := bitset.checkValidBit(n); err != nil {
		return false, err
	}
	idx := len(bitset.bitArray) - 1 - n/64
	return bitset.bitArray[idx]&(1<<(n%64)) >= 1, nil
}

// TestBits tests if multiple bit and returns a slice of bools that are true/false
// if the corresponding bits are set, and the number of set bits.
func (bitset *BitSet) TestBits(bits []int) ([]bool, int, error) {
	res := make([]bool, len(bitset.bitArray))
	numSet := 0
	for i, bit := range bits {
		isSet, err := bitset.Test(bit)
		if err != nil {
			return nil, 0, err
		}
		if isSet {
			numSet += 1
		}
		res[i] = isSet
	}
	return res, numSet, nil
}

// CountSetBits returns the number of set bits
func (bitset *BitSet) CountSetBits() int {
	sum := 0
	for _, word := range bitset.bitArray {
		sum += bits.OnesCount64(word)
	}
	return sum
}

// Or returns the result of bitset OR (|) other.
func (bitset *BitSet) Or(other *BitSet) *BitSet {
	smallerSet, greaterSet := bitset, other
	if bitset.size > other.size {
		smallerSet, greaterSet = other, bitset
	}
	newBitArray := make([]uint64, int(math.Ceil(float64(greaterSet.size)/64.0)))
	for i := len(smallerSet.bitArray) - 1; i >= 0; i-- {
		newBitArray[i] = smallerSet.bitArray[i] | greaterSet.bitArray[i]
	}
	return &BitSet{size: greaterSet.size, bitArray: newBitArray}
}

// And returns the result of bitset AND (&) other
func (bitset *BitSet) And(other *BitSet) *BitSet {
	smallerSet, greaterSet := bitset, other
	if bitset.size > other.size {
		smallerSet, greaterSet = other, bitset
	}
	newBitArray := make([]uint64, int(math.Ceil(float64(greaterSet.size)/64.0)))
	for i := len(smallerSet.bitArray) - 1; i >= 0; i-- {
		newBitArray[i] = smallerSet.bitArray[i] & greaterSet.bitArray[i]
	}
	return &BitSet{size: greaterSet.size, bitArray: newBitArray}
}

func (bitset *BitSet) String() string {
	buffer := bytes.NewBufferString("")
	for _, word := range bitset.bitArray {
		buffer.WriteString(fmt.Sprintf("%b", word))
	}
	return buffer.String()
}

func (bitset *BitSet) checkValidBit(n int) error {
	if n < 0 {
		return fmt.Errorf("test: n must be >= 0")
	}
	if n >= bitset.size {
		return fmt.Errorf("bit index %d out of range of bitset", n)
	}
	return nil
}
