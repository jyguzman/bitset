package bitset

import (
	"bytes"
	"fmt"
	"math/bits"
	"strings"
)

type BitSet struct {
	size  int // the number of bits the bitset holds
	words []uint64
}

// NewBitSetWithInitialSize initializes and returns a BitSet holding the given number of bits.
func NewBitSetWithInitialSize(numBits int) *BitSet {
	numWords := 1 + int(float64(numBits)/64.0)
	return &BitSet{
		size:  numBits,
		words: make([]uint64, numWords),
	}
}

// NewBitSet initializes and returns a BitSet with an initial size of 64.
func NewBitSet() *BitSet {
	return NewBitSetWithInitialSize(64)
}

// Size returns the number of bits the bitset holds
func (bs *BitSet) Size() int {
	return bs.size
}

// Set sets the Nth bit to 1.
func (bs *BitSet) Set(n int) {
	bs.resize(n)
	bs.set(n)
}

// SetBits sets multiple bits. This operation is atomic; if any bit is invalid,
// the bitset will roll back to its original state before attempting to set any of the
// bits.
func (bs *BitSet) SetBits(indices []int) {
	for _, idx := range indices {
		bs.Set(idx)
	}
}

// Clear zeroes the Nth bit. Errors if n < 0 or n >= bitset.size
func (bs *BitSet) Clear(n int) {
	bs.resize(n)
	bs.clear(n)
}

// ClearBits clears multiple bits. This operation is atomic; if any bit is invalid,
// the bitset will roll back to its original state before attempting to clear any of the
// bits.
func (bs *BitSet) ClearBits(indices []int) {
	for _, idx := range indices {
		bs.clear(idx)
	}
}

// ClearAll clears all bits.
func (bs *BitSet) ClearAll() {
	bs.words = make([]uint64, len(bs.words))
}

// Flip flips the Nth bit, i.e. 0 -> 1 or 1 -> 0.
func (bs *BitSet) Flip(n int) {
	bs.resize(n)
	bs.flip(n)
}

// FlipBits flips multiple bits. This operation is atomic; if any bit is invalid,
// the bitset will roll back to its original state before the attempt to flip the bits.
func (bs *BitSet) FlipBits(bits []int) {
	for _, idx := range bits {
		bs.Flip(idx)
	}
}

// Test checks if the Nth bit is set to 1. Errors if n < 0 or n >= bitset.size
func (bs *BitSet) Test(n int) bool {
	wordIdx, bitIdx := bs.getWordAndPos(n)
	return bs.words[wordIdx]&(1<<bitIdx) >= 1
}

// TestBits tests if multiple bits are set to 1. Returns a slice of bools that are true/false
// if the corresponding bits are set and the number of set bits.
func (bs *BitSet) TestBits(bits []int) ([]bool, int) {
	res, numSet := make([]bool, len(bits)), 0
	for i, bit := range bits {
		isSet := bs.Test(bit)
		if isSet {
			numSet += 1
		}
		res[i] = isSet
	}
	return res, numSet
}

// CountSetBits returns the number of set bits.
func (bs *BitSet) CountSetBits() int {
	count := 0
	for _, word := range bs.words {
		count += bits.OnesCount64(word)
	}
	return count
}

// Or sets the bits of the receiver to the result of the receiver OR (|) other.
func (bs *BitSet) Or(other *BitSet) {
	bitsLeft := bs.size
	for i, j := 0, 0; i < len(bs.words) && j < len(other.words); i, j = i+1, j+1 {
		bs.words[i] = mask(bs.words[i]|other.words[j], bitsLeft)
		bitsLeft -= 64
	}
}

// And sets the bits of the receiver to the result of the receiver AND (&) other.
func (bs *BitSet) And(other *BitSet) {
	bitsLeft := bs.size
	for i, j := 0, 0; i < len(bs.words) && j < len(other.words); i, j = i+1, j+1 {
		bs.words[i] = mask(bs.words[i]&other.words[j], bitsLeft)
		bitsLeft -= 64
	}
}

// Xor sets the bits of the receiver to the result of the receiver AND (&) other.
func (bs *BitSet) Xor(other *BitSet) {
	bitsLeft := bs.size
	for i, j := 0, 0; i < len(bs.words) && j < len(other.words); i, j = i+1, j+1 {
		bs.words[i] = mask(bs.words[i]^other.words[j], bitsLeft)
		bitsLeft -= 64
	}
}

// Not flips each bit of the bitset
func (bs *BitSet) Not() {
	bitsLeft := bs.size
	for i := range bs.words {
		bs.words[i] = mask(^bs.words[i], bitsLeft%64)
		bitsLeft -= 64
	}
}

// Any returns true if at least one bit is set
func (bs *BitSet) Any() bool {
	for _, word := range bs.words {
		if word != 0 {
			return true
		}
	}
	return false
}

// None returns true if no bits are set
func (bs *BitSet) None() bool {
	for _, word := range bs.words {
		if word != 0 {
			return false
		}
	}
	return true
}

// Or returns the result of bitset OR (|) other. The result's size will be equal to that of the
// larger bitset.
func Or(bs1 *BitSet, bs2 *BitSet) *BitSet {
	smallerSet, largerSet := bs1, bs2
	if bs1.size > bs2.size {
		smallerSet, largerSet = bs2, bs1
	}
	newBitArray := make([]uint64, len(largerSet.words))
	for i := len(smallerSet.words) - 1; i >= 0; i-- {
		newBitArray[i] = smallerSet.words[i] | largerSet.words[i]
	}
	return &BitSet{size: largerSet.size, words: newBitArray}
}

// And returns the result of bitset AND (&) other. The result's size will be equal to that of the
// larger bitset.
func And(bs1 *BitSet, bs2 *BitSet) *BitSet {
	smallerSet, largerSet := bs1, bs2
	if bs1.size > bs2.size {
		smallerSet, largerSet = bs2, bs1
	}
	newBitArray := make([]uint64, len(largerSet.words))
	for i := len(smallerSet.words) - 1; i >= 0; i-- {
		newBitArray[i] = smallerSet.words[i] & largerSet.words[i]
	}
	return &BitSet{size: largerSet.size, words: newBitArray}
}

// Not returns a new bitset obtained from flipping each bit of the input bitset.
func Not(bs *BitSet) *BitSet {
	newBitArray := make([]uint64, len(bs.words))
	for i := range bs.words {
		newBitArray[i] = ^bs.words[i]
	}
	return &BitSet{size: bs.size, words: newBitArray}
}

// Strings returns the representation of the bitset as a binary string.
func (bs *BitSet) String() string {
	buffer := bytes.Buffer{}
	for i := len(bs.words) - 1; i >= 0; i-- {
		buffer.WriteString(fmt.Sprintf("%.64b", bs.words[i]))
	}
	return strings.TrimLeft(buffer.String(), "0")
}

// set sets the Nth bit to 1.
func (bs *BitSet) set(n int) {
	wordIdx, bitIdx := bs.getWordAndPos(n)
	bs.words[wordIdx] |= 1 << bitIdx
}

// clear zeroes the Nth bit.
func (bs *BitSet) clear(n int) {
	wordIdx, bitIdx := bs.getWordAndPos(n)
	bs.words[wordIdx] &= ^(1 << bitIdx)
}

// flip flips the Nth bit, i.e. 0 -> 1 or 1 -> 0.
func (bs *BitSet) flip(n int) {
	wordIdx, bitIdx := bs.getWordAndPos(n)
	bs.words[wordIdx] ^= 1 << bitIdx
}

func (bs *BitSet) getWordAndPos(n int) (int, int) {
	return n / 64, n % 64
}

func (bs *BitSet) resize(newSize int) {
	if newSize >= bs.size {
		bs.size = newSize
		bs.words = append(bs.words, make([]uint64, len(bs.words))...)
	}
}

func (bs *BitSet) checkValidBit(n int) error {
	if n < 0 {
		return fmt.Errorf("test: n must be >= 0")
	}
	if n >= bs.size {
		return fmt.Errorf("bit index %d out of range of bitset of size %d", n, bs.size)
	}
	return nil
}

// mask retains the first n bits of a word and zeroes out the rest, returning the result.
// If n is invalid the original word is returned.
func mask(word uint64, n int) uint64 {
	if n <= 0 || n >= 64 {
		return word
	}
	return word & ((1 << n) - 1)
}
