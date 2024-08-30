package bitset

import (
	"bytes"
	"fmt"
	"math"
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
func (bitset *BitSet) Size() int {
	return bitset.size
}

// Set sets the Nth bit to 1. Errors if n < 0 or n >= bitset.size
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

		// out of bounds: roll back by un-setting bits up to the invalid bit
		if err != nil {
			for j, ob := range originalBits {
				// if the bit was 0 before being set, clear it back to 0
				if ob == 0 {
					bitset.clear(indices[j])
				}
			}
			return err
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

		// out of bounds: roll back by un-clearing bits up to invalid bit
		if err != nil {
			for j, ob := range originalBits {
				// if the bit was 1 before being cleared, set it back to 1
				if ob == 1 {
					bitset.set(indices[j])
				}
			}
			return err
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
	bitset.words = make([]uint64, int(math.Ceil(float64(bitset.size)/64.0)))
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
// the bitset will roll back to its original state before the attempt to flip the bits.
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

// Test checks if the Nth bit is set to 1. Errors if n < 0 or n >= bitset.size
func (bitset *BitSet) Test(n int) (bool, error) {
	if err := bitset.checkValidBit(n); err != nil {
		return false, err
	}
	return bitset.test(n), nil
}

// TestBits tests if multiple bits are set to 1. Returns a slice of bools that are true/false
// if the corresponding bits are set and the number of set bits.
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

// CountSetBits returns the number of set bits.
func (bitset *BitSet) CountSetBits() int {
	count := 0
	for _, word := range bitset.words {
		count += bits.OnesCount64(word)
	}
	return count
}

// Or sets the bits of the receiver to the result of the receiver OR (|) other.
func (bitset *BitSet) Or(other *BitSet) {
	//recvBitsLeft, otherBitsLeft := bitset.size, other.size
	recvBitsLeft := bitset.size
	for i, j := len(bitset.words)-1, len(other.words)-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		btWord, otherWord := bitset.words[i], other.words[j]
		//if recvBitsLeft < 64 || otherBitsLeft < 64 {
		//	maskLen := int(math.Min(float64(recvBitsLeft), float64(otherBitsLeft)))
		//	bitset.words[i] = mask(btWord|otherWord, maskLen)
		//} else {
		//	bitset.words[i] = btWord | otherWord
		//}
		if recvBitsLeft < 64 {
			//maskLen := int(math.Min(float64(recvBitsLeft), float64(otherBitsLeft)))
			bitset.words[i] = mask(btWord|otherWord, recvBitsLeft)
		} else {
			bitset.words[i] = btWord | otherWord
		}
		recvBitsLeft -= 64
		//otherBitsLeft -= 64
	}
}

// And sets the bits of the receiver to the result of the receiver AND (&) other.
func (bitset *BitSet) And(other *BitSet) {
	recvBitsLeft, otherBitsLeft := bitset.size, other.size
	for i, j := len(bitset.words)-1, len(other.words)-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		btWord, otherWord := bitset.words[i], other.words[j]
		if recvBitsLeft < 64 || otherBitsLeft < 64 {
			maskLen := int(math.Min(float64(recvBitsLeft), float64(otherBitsLeft)))
			bitset.words[i] = mask(btWord&otherWord, maskLen)
		} else {
			bitset.words[i] = btWord & otherWord
		}
		recvBitsLeft -= 64
		otherBitsLeft -= 64
	}
}

// Xor sets the bits of the receiver to the result of the receiver AND (&) other.
func (bitset *BitSet) Xor(other *BitSet) {
	recvBitsLeft, otherBitsLeft := bitset.size, other.size
	for i, j := len(bitset.words)-1, len(other.words)-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		btWord, otherWord := bitset.words[i], other.words[j]
		if recvBitsLeft < 64 || otherBitsLeft < 64 {
			maskLen := int(math.Min(float64(recvBitsLeft), float64(otherBitsLeft)))
			bitset.words[i] = mask(btWord^otherWord, maskLen)
		} else {
			bitset.words[i] = btWord ^ otherWord
		}
		recvBitsLeft -= 64
		otherBitsLeft -= 64
	}
}

// Not flips each bit of the bitset
func (bitset *BitSet) Not() {
	bitsLeft := bitset.size
	for i := range bitset.words {
		if bitsLeft < 64 {
			bitset.words[i] = mask(^bitset.words[i], bitsLeft)
		} else {
			bitset.words[i] = ^bitset.words[i]
		}
		bitsLeft -= 64
	}
}

// Or returns the result of bitset OR (|) other.
func Or(bs1 *BitSet, bs2 *BitSet) *BitSet {
	smallerSet, largerSet := bs1, bs2
	if bs1.size > bs2.size {
		smallerSet, largerSet = bs2, bs1
	}
	newBitArray := make([]uint64, int(math.Ceil(float64(largerSet.size)/64.0)))
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
	newBitArray := make([]uint64, int(math.Ceil(float64(largerSet.size)/64.0)))
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

func (bitset *BitSet) String() string {
	buffer := bytes.Buffer{}
	for i, word := range bitset.words {
		if word == 0 && i != len(bitset.words)-1 {
			continue
		}
		buffer.WriteString(fmt.Sprintf("%b", word))
	}
	return buffer.String()
}

// set sets the Nth bit to 1.
func (bitset *BitSet) set(n int) {
	wordIdx, bitIdx := bitset.getWordAndPos(n)
	bitset.words[wordIdx] |= 1 << bitIdx
}

// clear zeroes the Nth bit.
func (bitset *BitSet) clear(n int) {
	wordIdx, bitIdx := bitset.getWordAndPos(n)
	bitset.words[wordIdx] &= ^(1 << bitIdx)
}

// flip flips the Nth bit, i.e. 0 -> 1 or 1 -> 0.
func (bitset *BitSet) flip(n int) {
	wordIdx, bitIdx := bitset.getWordAndPos(n)
	bitset.words[wordIdx] ^= 1 << bitIdx
}

// test checks if the Nth bit is set.
func (bitset *BitSet) test(n int) bool {
	wordIdx, bitIdx := bitset.getWordAndPos(n)
	//fmt.Println("wordIdx:", wordIdx)
	return bitset.words[wordIdx]&(1<<bitIdx) >= 1
}

func (bitset *BitSet) getWordAndPos(n int) (int, int) {
	fmt.Println("arr:", bitset.words, bitset.size)
	fmt.Println("n:", n)
	bitIdx := n % 64
	return len(bitset.words) - n/64 - 1, bitIdx
}

// mask retains the first n bits of a word and zeroes out the rest, returning the result
func mask(word uint64, n int) uint64 {
	return word & ((1 << n) - 1)
}

func (bitset *BitSet) checkValidBit(n int) error {
	if n < 0 {
		return fmt.Errorf("test: n must be >= 0")
	}
	if n >= bitset.size {
		return fmt.Errorf("bit index %d out of range of bitset of size %d", n, bitset.size)
	}
	return nil
}
