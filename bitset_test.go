package bitset

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"testing"
)

func TestBitSet_BitArrayLenDivisibleBy64(t *testing.T) {
	numBits := 512
	bs := NewBitSet(numBits)
	want := 8
	if len(bs.bitArray) != want {
		t.Errorf("BitSet has length %d, want %d", len(bs.bitArray), want)
	}
}

func TestBitSet_BitArrayLenNotDivisibleBy64(t *testing.T) {
	numBits := 512
	bs := NewBitSet(numBits + 1)
	want := 9
	if len(bs.bitArray) != want {
		t.Errorf("BitSet has length %d, want %d", len(bs.bitArray), want)
	}
	bs = NewBitSet(numBits - 1)
	want = 8
	if len(bs.bitArray) != want {
		t.Errorf("BitSet has length %d, want %d", len(bs.bitArray), want)
	}
}

func TestBitSet_Test(t *testing.T) {
	bitArray := []uint64{uint64(math.Pow(2.0, 63.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 1, so bits 0 and 63 should be set
	bs := &BitSet{size: 64, bitArray: bitArray}

	isSet, err := bs.Test(63)
	if !isSet {
		t.Errorf("BitSet.Test(63) == false, want true")
	}

	isSet, err = bs.Test(0)
	if !isSet {
		t.Errorf("BitSet.Test(0) == false, want true")
	}

	isSet, err = bs.Test(30)
	if isSet {
		t.Errorf("BitSet.Test(30) == true, want false")
	}

	err = bs.Set(64)
	if err == nil {
		t.Errorf("BitSet Test(): should have returned out of bounds error")
	}

	err = bs.Set(-1)
	if err == nil {
		t.Errorf("BitSet Test(): should have returned out of bounds error")
	}
}

func TestBitSet_TestBits(t *testing.T) {
	bitArray := []uint64{uint64(math.Pow(2.0, 63.0)) + uint64(math.Pow(2.0, 30.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 2^30 + 1, so bits 0, 30, and 63 should be set
	bs := &BitSet{size: 64, bitArray: bitArray}

	bitsToTest := []int{0, 15, 30, 45, 63}

	bools, numSet, err := bs.TestBits(bitsToTest)

	if err != nil {
		t.Errorf("BitSet.TestBits() failed: %v", err)
	}

	if numSet != 3 {
		t.Errorf("BitSet.TestBits() want 3, got %d", numSet)
	}

	if !slices.Equal(bools, []bool{true, false, true, false, true}) {
		t.Errorf("BitSet.TestBits() want [true, false, true, false, true], got %v", bools)
	}
}

func TestBitSet_Set(t *testing.T) {
	bs := NewBitSet(64)

	err := bs.Set(64)
	if err == nil {
		t.Errorf("BitSet should have returned out of bounds error for invalid bit 64")
	}
	err = bs.Set(-1)
	if err == nil {
		t.Errorf("BitSet should have returned out of bounds error for invalid bit -1")
	}

	err = bs.Set(0)
	isSet, err := bs.Test(0)
	if !isSet {
		t.Errorf("BitSet.Test(0) == false, want true")
	}

	err = bs.Set(63)
	isSet, err = bs.Test(63)
	if !isSet {
		t.Errorf("BitSet.Test(63) == false, want true")
	}

	err = bs.Set(0)
	isSet, err = bs.Test(0)
	if !isSet {
		t.Errorf("BitSet.Test(0) on already set bit == false, want true")
	}
}

func TestBitSet_SetBits(t *testing.T) {
	bs := NewBitSet(64)
	bitsToSet := []int{0, 63, 0, 5, 10}
	if err := bs.SetBits(bitsToSet); err != nil {
		t.Errorf("BitSet.SetBits() failed: %v", err)
	}
	bools, numSet, err := bs.TestBits(bitsToSet)
	if err != nil {
		t.Errorf("BitSet.TestBits() failed in Test_SetBits: %v", err)
	}
	if numSet != len(bitsToSet) {
		t.Errorf("BitSet.TestBits() want %d, got %d", len(bitsToSet), numSet)
	}
	if !slices.Equal(bools, []bool{true, true, true, true, true}) {
		t.Errorf("BitSet.TestBits() want [true, false, true, false, true], got %v", bools)
	}
}

func TestBitSet_Clear(t *testing.T) {
	bitArray := []uint64{uint64(math.Pow(2.0, 63.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 1, so bits 0 and 63 should be set
	bs := &BitSet{size: 64, bitArray: bitArray}

	err := bs.Clear(64)
	if err == nil {
		t.Errorf("BitSet should have returned out of bounds error for invalid bit 64")
	}
	err = bs.Clear(-1)
	if err == nil {
		t.Errorf("BitSet should have returned out of bounds error for invalid bit -1")
	}

	err = bs.Clear(0)
	isSet, err := bs.Test(0)
	if isSet {
		t.Errorf("BitSet.Test(0) == true, want false")
	}

	err = bs.Clear(63)
	isSet, err = bs.Test(63)
	if isSet {
		t.Errorf("BitSet.Test(63) == true, want false")
	}

	err = bs.Clear(30)
	isSet, err = bs.Test(30)
	if isSet {
		t.Errorf("BitSet.Test(0) on zero bit == true, want false")
	}

	err = bs.Clear(0)
	isSet, err = bs.Test(0)
	if isSet {
		t.Errorf("BitSet.Test(0) on cleared bit == true, want false")
	}
}

func TestBitSet_ClearBits(t *testing.T) {
	bitArray := []uint64{uint64(math.Pow(2.0, 63.0)) + uint64(math.Pow(2.0, 30.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 2^30 + 1, so bits 0, 30, and 63 should be set
	bs := &BitSet{size: 64, bitArray: bitArray}

	bitsToClear := []int{0, 15, 30, 45, 63}
	if err := bs.ClearBits(bitsToClear); err != nil {
		t.Errorf("BitSet.ClearBits() failed: %v", err)
	}
	bools, numSet, err := bs.TestBits(bitsToClear)
	if err != nil {
		t.Errorf("BitSet.TestBits() failed in Test_ClearBits: %v", err)
	}
	if numSet > 0 {
		t.Errorf("BitSet.TestBits() want %d, got %d", 0, numSet)
	}
	if !slices.Equal(bools, []bool{false, false, false, false, false}) {
		t.Errorf("BitSet.TestBits() in ClearBits: want [true, false, true, false, true], got %v", bools)
	}
}

func TestBitSet_Flip(t *testing.T) {
	bitArray := []uint64{uint64(math.Pow(2.0, 63.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 1, so bits 0 and 63 should be set
	bs := &BitSet{size: 64, bitArray: bitArray}

	err := bs.Flip(64)
	if err == nil {
		t.Errorf("BitSet.Flip() should have returned out of bounds error for invalid bit 64")
	}
	err = bs.Flip(-1)
	if err == nil {
		t.Errorf("BitSet.Flip() should have returned out of bounds error for invalid bit -1")
	}

	err = bs.Flip(0)
	isSet, err := bs.Test(0)
	if isSet {
		t.Errorf("BitSet.Flip(0) on 1 bit should be 0, still have 1")
	}

	err = bs.Flip(0)
	isSet, err = bs.Test(0)
	if !isSet {
		t.Errorf("BitSet.Flip(0) on 0 bit should be 1, still have 0")
	}

	err = bs.Flip(63)
	isSet, err = bs.Test(63)
	if isSet {
		t.Errorf("BitSet.Flip(0) on 1 bit should be 0, still have 1")
	}

	err = bs.Flip(63)
	isSet, err = bs.Test(63)
	if !isSet {
		t.Errorf("BitSet.Flip(0) on 0 bit should be 1, still have 0")
	}

	err = bs.Flip(30)
	isSet, err = bs.Test(30)
	if !isSet {
		t.Errorf("BitSet.Test(0) on 0 should be 1, still have 0")
	}
}

func TestBitSet_FlipBits(t *testing.T) {
	bitArray := []uint64{uint64(math.Pow(2.0, 63.0)) + uint64(math.Pow(2.0, 30.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 2^30 + 1, so bits 0, 30, and 63 should be set
	bs := &BitSet{size: 64, bitArray: bitArray}

	bitsToFlip := []int{0, 15, 30, 45, 63}
	if err := bs.FlipBits(bitsToFlip); err != nil {
		t.Errorf("BitSet.ClearBits() failed: %v", err)
	}
	bools, numSet, err := bs.TestBits(bitsToFlip)
	if err != nil {
		t.Errorf("BitSet.TestBits() failed in Test_ClearBits: %v", err)
	}
	if numSet != 2 {
		t.Errorf("BitSet.TestBits() want %d, got %d", 2, numSet)
	}
	if !slices.Equal(bools, []bool{false, true, false, true, false}) {
		t.Errorf("BitSet.TestBits() in ClearBits: want [true, false, true, false, true], got %v", bools)
	}
}

func TestBitSet_Or(t *testing.T) {
	a := NewBitSet(20)
	b := NewBitSet(20)

	aToSet := []int{0, 1, 2, 4}
	if err := a.SetBits(aToSet); err != nil {
		t.Error(err)
	}
	res := a.Or(b)
	fmt.Println(res)
}

func TestBitSet_And(t *testing.T) {
	a := NewBitSet(20)
	b := NewBitSet(13)

	aToSet := []int{1, 5, 10, 15, 17}
	bToSet := []int{0, 3, 6, 9, 12}
	if err := a.SetBits(aToSet); err != nil {
		t.Error(err)
	}
	if err := b.SetBits(bToSet); err != nil {
		t.Error(err)
	}
	fmt.Println("a:", a, "b: ", b)
	res := a.And(b)
	fmt.Printf("%b\n", 0b00101000010000100010&0b1001001001001)
	fmt.Println(res)
}

func TestBitSet_Not(t *testing.T) {
	a := NewBitSet(20)

	aToSet := []int{1, 5, 10, 15, 17}
	if err := a.SetBits(aToSet); err != nil {
		t.Error(err)
	}
	fmt.Println(a)
	a.Not()
	fmt.Println(a)
}

func TestBitSet_String(t *testing.T) {
	numBits := rand.Intn(7)
	numBitsToSet := rand.Intn(numBits)
	bits, setBits := make([]int, numBitsToSet), make(map[int]bool)
	for i := 0; i < numBitsToSet; i++ {
		bit := rand.Intn(numBits)
		setBits[bit], bits[i] = true, bit
	}
	bs := NewBitSet(numBits)
	if err := bs.SetBits(bits); err != nil {
		t.Error(err)
	}
	str := bs.String()
	fmt.Println()
	for i := len(str) - 1; i >= 0; i-- {
		_, ok := setBits[i]
		if str[i] == '1' && !ok {
			t.Errorf("SetBits: failed for bit %d, is %d but want %d", i, int(str[i]-'0'), 0)
		}
		if str[i] == '0' && ok {
			t.Errorf("SetBits: failed for bit %d, is %d but want %d", i, int(str[i]-'0'), 1)
		}
	}
}

func TestBitSet_Count(t *testing.T) {
	bs := NewBitSet(0)
	count, want := bs.CountSetBits(), 0
	if count != want {
		t.Errorf("BitSet.CountSetBits() with empty bitset: got %d, want %d", count, want)
	}
	numBits := 512
	bs, setBits := NewBitSet(numBits), make(map[uint64]bool)
	numBitsToSet := rand.Intn(numBits)
	bits := make([]int, numBitsToSet)
	for i := 0; i < numBitsToSet; i++ {
		bit := rand.Intn(numBits)
		setBits[uint64(bit)], bits[i] = true, bit
	}
	if err := bs.SetBits(bits); err != nil {
		t.Error(err)
	}
	count, want = bs.CountSetBits(), len(setBits)
	if count != want {
		t.Errorf("BitSet has count %d, want %d", count, want)
	}
}

func Test_Do(t *testing.T) {
	fmt.Printf("%b\n", 0b00000|0b1001)
}
