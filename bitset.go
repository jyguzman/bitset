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
	bitset.set(n)
	return nil
}

// SetBits sets multiple bits. This operation is atomic; if any bit is invalid,
// the bitset will roll back to its original state before attempting to set any of the
// bits.
func (bitset *BitSet) SetBits(indices []int) error {
	var originalBits []int8

	for _, idx := range indices {
		wasSet, err := bitset.Test(idx)

		// out of bounds: roll back
		if err != nil {
			for j, ob := range originalBits {
				// if the bit was 0 before being set, clear it back to 0
				if ob == 0 {
					bitset.clear(indices[j])
				}
			}
		}

		if wasSet {
			originalBits = append(originalBits, 1)
		} else {
			originalBits = append(originalBits, 0)
		}

		bitset.set(idx)
	}

	return nil
}

// Clear zeroes the Nth bit. Errors if n < 0 or n >= bitset.size
func (bitset *BitSet) Clear(n int) error {
	if err := bitset.checkValidBit(n); err != nil {
		return err
	}
	bitset.clear(n)
	return nil
}

// ClearBits clears multiple bits. This operation is atomic; if any bit is invalid,
// the bitset will roll back to its original state before attempting to clear any of the
// bits.
func (bitset *BitSet) ClearBits(indices []int) error {
	var originalBits []int8

	for _, idx := range indices {
		wasSet, err := bitset.Test(idx)

		// out of bounds: roll back
		if err != nil {
			for j, ob := range originalBits {
				// if the bit was 1 before being cleared, set it back to 1
				if ob == 1 {
					bitset.set(indices[j])
				}
			}
		}

		if wasSet {
			originalBits = append(originalBits, 1)
		} else {
			originalBits = append(originalBits, 0)
		}

		bitset.clear(idx)
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
	bitset.flip(n)
	return nil
}

// FlipBits flips multiple bits. This operation is atomic; if any bit is invalid,
// the bitset will roll back to its original state.
func (bitset *BitSet) FlipBits(bits []int) error {
	for _, idx := range bits {
		err := bitset.Flip(idx)
		if err != nil {
			for _, i := range bits {
				if i == idx {
					return err
				}
				bitset.flip(i)
			}
		}
	}
	return nil
}

// Test checks if the Nth bit is set. Errors if n < 0 or n >= bitset.size
func (bitset *BitSet) Test(n int) (bool, error) {
	if err := bitset.checkValidBit(n); err != nil {
		return false, err
	}
	return bitset.test(n), nil
}

// TestBits tests if multiple bit and returns a slice of bools that are true/false
// if the corresponding bits are set, and the number of set bits.
func (bitset *BitSet) TestBits(bits []int) ([]bool, int, error) {
	res, numSet := make([]bool, len(bits)), 0
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
	count := 0
	for _, word := range bitset.bitArray {
		count += bits.OnesCount64(word)
	}
	return count
}

// Or returns the result of bitset OR (|) other.
func (bitset *BitSet) Or(other *BitSet) *BitSet {
	smallerSet, largerSet := bitset, other
	if bitset.size > other.size {
		smallerSet, largerSet = other, bitset
	}
	newBitArray := make([]uint64, int(math.Ceil(float64(largerSet.size)/64.0)))
	for i := len(smallerSet.bitArray) - 1; i >= 0; i-- {
		newBitArray[i] = smallerSet.bitArray[i] | largerSet.bitArray[i]
	}
	return &BitSet{size: largerSet.size, bitArray: newBitArray}
}

// And returns the result of bitset AND (&) other
func (bitset *BitSet) And(other *BitSet) *BitSet {
	smallerSet, largerSet := bitset, other
	if bitset.size > other.size {
		smallerSet, largerSet = other, bitset
	}
	newBitArray := make([]uint64, int(math.Ceil(float64(largerSet.size)/64.0)))
	for i := len(smallerSet.bitArray) - 1; i >= 0; i-- {
		newBitArray[i] = smallerSet.bitArray[i] & largerSet.bitArray[i]
	}
	return &BitSet{size: largerSet.size, bitArray: newBitArray}
}

// Not flips each bit of the bitset
func (bitset *BitSet) Not() {
	for i := range bitset.bitArray {
		bitset.bitArray[i] = ^bitset.bitArray[i]
	}
}

func (bitset *BitSet) String() string {
	buffer := bytes.NewBufferString("")
	for _, word := range bitset.bitArray {
		buffer.WriteString(fmt.Sprintf("%b", word))
	}
	return buffer.String()
}

// set sets the Nth bit to 1.
func (bitset *BitSet) set(n int) {
	word, idx := bitset.getWordAndPos(n)
	bitset.bitArray[word] ^= 1 << idx
}

// clear zeroes the Nth bit.
func (bitset *BitSet) clear(n int) {
	word, idx := bitset.getWordAndPos(n)
	bitset.bitArray[word] &= ^(1 << idx)
}

// flip flips the Nth bit, i.e. 0 -> 1 or 1 -> 0.
func (bitset *BitSet) flip(n int) {
	word, idx := bitset.getWordAndPos(n)
	bitset.bitArray[word] |= 1 << idx
}

// test checks if the Nth bit is set.
func (bitset *BitSet) test(n int) bool {
	word, idx := bitset.getWordAndPos(n)
	return bitset.bitArray[word]&(1<<idx) >= 1
}

func (bitset *BitSet) getWordAndPos(n int) (int, int) {
	return len(bitset.bitArray) - 1 - n/64, n % 64
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
