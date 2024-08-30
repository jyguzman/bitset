package bitset

import (
	"bytes"
	"fmt"
	"math/bits"
)

type BitSet struct {
	size  int // the number of bits the bitset holds
	words []uint64
}

// NewBitSet initializes and returns a BitSet holding the given number of bits
func NewBitSet(numBits int) *BitSet {
	numWords := 1 + int(float64(numBits)/64.0)
	return &BitSet{
		size:  numBits,
		words: make([]uint64, numWords),
	}
}

// Size returns the number of bits the bitset holds
func (bs *BitSet) Size() int {
	return bs.size
}

// Set sets the Nth bit to 1. Errors if n < 0 or n >= bitset.size
func (bs *BitSet) Set(n int) error {
	if err := bs.checkValidBit(n); err != nil {
		return err
	}
	bs.set(n)
	return nil
}

// SetBits sets multiple bits. This operation is atomic; if any bit is invalid,
// the bitset will roll back to its original state before attempting to set any of the
// bits.
func (bs *BitSet) SetBits(indices []int) error {
	var originalBits []int8
	for _, idx := range indices {
		wasSet, err := bs.Test(idx)

		// out of bounds: roll back by un-setting bits up to the invalid bit
		if err != nil {
			for j, ob := range originalBits {
				// if the bit was 0 before being set, clear it back to 0
				if ob == 0 {
					bs.clear(indices[j])
				}
			}
			return err
		}

		if wasSet {
			originalBits = append(originalBits, 1)
		} else {
			originalBits = append(originalBits, 0)
		}

		bs.set(idx)
	}

	return nil
}

// Clear zeroes the Nth bit. Errors if n < 0 or n >= bitset.size
func (bs *BitSet) Clear(n int) error {
	if err := bs.checkValidBit(n); err != nil {
		return err
	}
	bs.clear(n)
	return nil
}

// ClearBits clears multiple bits. This operation is atomic; if any bit is invalid,
// the bitset will roll back to its original state before attempting to clear any of the
// bits.
func (bs *BitSet) ClearBits(indices []int) error {
	var originalBits []int8

	for _, idx := range indices {
		wasSet, err := bs.Test(idx)

		// out of bounds: roll back by un-clearing bits up to invalid bit
		if err != nil {
			for j, ob := range originalBits {
				// if the bit was 1 before being cleared, set it back to 1
				if ob == 1 {
					bs.set(indices[j])
				}
			}
			return err
		}

		if wasSet {
			originalBits = append(originalBits, 1)
		} else {
			originalBits = append(originalBits, 0)
		}

		bs.clear(idx)
	}

	return nil
}

// ClearAll clears all bits.
func (bs *BitSet) ClearAll() {
	bs.words = make([]uint64, len(bs.words))
}

// Flip flips the Nth bit, i.e. 0 -> 1 or 1 -> 0. Errors if n < 0 or n >= bitset.size
func (bs *BitSet) Flip(n int) error {
	if err := bs.checkValidBit(n); err != nil {
		return err
	}
	bs.flip(n)
	return nil
}

// FlipBits flips multiple bits. This operation is atomic; if any bit is invalid,
// the bitset will roll back to its original state before the attempt to flip the bits.
func (bs *BitSet) FlipBits(bits []int) error {
	for _, idx := range bits {
		err := bs.Flip(idx)
		if err != nil {
			for _, i := range bits {
				if i == idx {
					return err
				}
				bs.flip(i)
			}
		}
	}
	return nil
}

// Test checks if the Nth bit is set to 1. Errors if n < 0 or n >= bitset.size
func (bs *BitSet) Test(n int) (bool, error) {
	if err := bs.checkValidBit(n); err != nil {
		return false, err
	}
	return bs.test(n), nil
}

// TestBits tests if multiple bits are set to 1. Returns a slice of bools that are true/false
// if the corresponding bits are set and the number of set bits.
func (bs *BitSet) TestBits(bits []int) ([]bool, int, error) {
	res, numSet := make([]bool, len(bits)), 0
	for i, bit := range bits {
		isSet, err := bs.Test(bit)
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
		bs.words[i] = mask(bs.words[i]|other.words[j], bitsLeft%64)
		bitsLeft -= 64
	}
}

// And sets the bits of the receiver to the result of the receiver AND (&) other.
func (bs *BitSet) And(other *BitSet) {
	bitsLeft := bs.size
	for i, j := 0, 0; i < len(bs.words) && j < len(other.words); i, j = i+1, j+1 {
		bs.words[i] = mask(bs.words[i]&other.words[j], bitsLeft%64)
		bitsLeft -= 64
	}
}

// Xor sets the bits of the receiver to the result of the receiver AND (&) other.
func (bs *BitSet) Xor(other *BitSet) {
	bitsLeft := bs.size
	for i, j := 0, 0; i < len(bs.words) && j < len(other.words); i, j = i+1, j+1 {
		bs.words[i] = mask(bs.words[i]^other.words[j], bitsLeft%64)
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

// Or returns the result of bitset OR (|) other.
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

// And returns the result of bitset AND (&) other
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

// Not returns a new bitset obtained from flipping each bit of the input bitset
func Not(bs *BitSet) *BitSet {
	newBitArray := make([]uint64, len(bs.words))
	for i := range bs.words {
		newBitArray[i] = ^bs.words[i]
	}
	return &BitSet{size: bs.size, words: newBitArray}
}

func (bs *BitSet) String() string {
	buffer := bytes.Buffer{}
	for i := len(bs.words) - 1; i >= 0; i-- {
		word := bs.words[i]
		if word == 0 {
			continue
		}
		if i+1 < len(bs.words) && bs.words[i+1] != 0 {
			buffer.WriteString(fmt.Sprintf("%.64b", word))
		} else {
			buffer.WriteString(fmt.Sprintf("%b", word))
		}
	}
	return buffer.String()
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

// test checks if the Nth bit is set.
func (bs *BitSet) test(n int) bool {
	wordIdx, bitIdx := bs.getWordAndPos(n)
	return bs.words[wordIdx]&(1<<bitIdx) >= 1
}

func (bs *BitSet) getWordAndPos(n int) (int, int) {
	return n / 64, n % 64
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
// If n == 0 the original word is returned
func mask(word uint64, n int) uint64 {
	if n == 0 {
		return word
	}
	return word & ((1 << n) - 1)
}
