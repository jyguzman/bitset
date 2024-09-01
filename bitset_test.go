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
	bs := NewBitSetInitialSize(numBits)
	want := 8
	if len(bs.words) != want {
		t.Errorf("BitSet has length %d, want %d", len(bs.words), want)
	}
}

func TestBitSet_BitArrayLenNotDivisibleBy64(t *testing.T) {
	numBits := 512
	bs := NewBitSetInitialSize(numBits + 1)
	want := 9
	if len(bs.words) != want {
		t.Errorf("BitSet has length %d, want %d", len(bs.words), want)
	}
	bs = NewBitSetInitialSize(numBits - 1)
	want = 8
	if len(bs.words) != want {
		t.Errorf("BitSet has length %d, want %d", len(bs.words), want)
	}
}

func TestBitSet_Test(t *testing.T) {
	words := []uint64{uint64(math.Pow(2.0, 63.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 1, so bits 0 and 63 should be set
	bs := &BitSet{size: 64, words: words}

	isSet := bs.Test(63)
	if !isSet {
		t.Errorf("BitSet.Test(63) == false, want true")
	}

	isSet = bs.Test(0)
	if !isSet {
		t.Errorf("BitSet.Test(0) == false, want true")
	}

	isSet = bs.Test(30)
	if isSet {
		t.Errorf("BitSet.Test(30) == true, want false")
	}

	bs.Set(64)

	bs.Set(-1)
}

func TestBitSet_TestBits(t *testing.T) {
	words := []uint64{uint64(math.Pow(2.0, 63.0)) + uint64(math.Pow(2.0, 30.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 2^30 + 1, so bits 0, 30, and 63 should be set
	bs := &BitSet{size: 64, words: words}

	bitsToTest := []int{0, 15, 30, 45, 63}

	bools, numSet := bs.TestBits(bitsToTest)

	if numSet != 3 {
		t.Errorf("BitSet.TestBits() want 3, got %d", numSet)
	}

	if !slices.Equal(bools, []bool{true, false, true, false, true}) {
		t.Errorf("BitSet.TestBits() want [true, false, true, false, true], got %v", bools)
	}
}

func TestBitSet_Set(t *testing.T) {
	bs := NewBitSetInitialSize(64)

	bs.Set(64)
	bs.Set(-1)

	bs.Set(0)
	isSet := bs.Test(0)
	if !isSet {
		t.Errorf("BitSet.Test(0) == false, want true")
	}

	bs.Set(63)
	isSet = bs.Test(63)
	if !isSet {
		t.Errorf("BitSet.Test(63) == false, want true")
	}

	bs.Set(0)
	isSet = bs.Test(0)
	if !isSet {
		t.Errorf("BitSet.Test(0) on already set bit == false, want true")
	}
}

func TestBitSet_SetBits(t *testing.T) {
	bs := NewBitSetInitialSize(64)
	bitsToSet := []int{0, 63, 0, 5, 10}
	bs.SetBits(bitsToSet)
	bools, numSet := bs.TestBits(bitsToSet)

	if numSet != len(bitsToSet) {
		t.Errorf("BitSet.TestBits() want %d, got %d", len(bitsToSet), numSet)
	}
	if !slices.Equal(bools, []bool{true, true, true, true, true}) {
		t.Errorf("BitSet.TestBits() want [true, false, true, false, true], got %v", bools)
	}
}

func TestBitSet_Clear(t *testing.T) {
	words := []uint64{uint64(math.Pow(2.0, 63.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 1, so bits 0 and 63 should be set
	bs := &BitSet{size: 64, words: words}

	bs.Clear(64)

	bs.Clear(-1)

	bs.Clear(0)
	isSet := bs.Test(0)
	if isSet {
		t.Errorf("BitSet.Test(0) == true, want false")
	}

	bs.Clear(63)
	isSet = bs.Test(63)
	if isSet {
		t.Errorf("BitSet.Test(63) == true, want false")
	}

	bs.Clear(30)
	isSet = bs.Test(30)
	if isSet {
		t.Errorf("BitSet.Test(0) on zero bit == true, want false")
	}

	bs.Clear(0)
	isSet = bs.Test(0)
	if isSet {
		t.Errorf("BitSet.Test(0) on cleared bit == true, want false")
	}
}

func TestBitSet_ClearBits(t *testing.T) {
	words := []uint64{uint64(math.Pow(2.0, 63.0)) + uint64(math.Pow(2.0, 30.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 2^30 + 1, so bits 0, 30, and 63 should be set
	bs := &BitSet{size: 64, words: words}

	bitsToClear := []int{0, 15, 30, 45, 63}
	bs.ClearBits(bitsToClear)
	bools, numSet := bs.TestBits(bitsToClear)

	if numSet > 0 {
		t.Errorf("BitSet.TestBits() want %d, got %d", 0, numSet)
	}
	if !slices.Equal(bools, []bool{false, false, false, false, false}) {
		t.Errorf("BitSet.TestBits() in ClearBits: want [true, false, true, false, true], got %v", bools)
	}
}

func TestBitSet_Flip(t *testing.T) {
	words := []uint64{uint64(math.Pow(2.0, 63.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 1, so bits 0 and 63 should be set
	bs := &BitSet{size: 64, words: words}

	bs.Flip(64)

	bs.Flip(-1)

	bs.Flip(0)
	isSet := bs.Test(0)
	if isSet {
		t.Errorf("BitSet.Flip(0) on 1 bit should be 0, still have 1")
	}

	bs.Flip(0)
	isSet = bs.Test(0)
	if !isSet {
		t.Errorf("BitSet.Flip(0) on 0 bit should be 1, still have 0")
	}

	bs.Flip(63)
	isSet = bs.Test(63)
	if isSet {
		t.Errorf("BitSet.Flip(0) on 1 bit should be 0, still have 1")
	}

	bs.Flip(63)
	isSet = bs.Test(63)
	if !isSet {
		t.Errorf("BitSet.Flip(0) on 0 bit should be 1, still have 0")
	}

	bs.Flip(30)
	isSet = bs.Test(30)
	if !isSet {
		t.Errorf("BitSet.Test(0) on 0 should be 1, still have 0")
	}
}

func TestBitSet_FlipBits(t *testing.T) {
	words := []uint64{uint64(math.Pow(2.0, 63.0)) + uint64(math.Pow(2.0, 30.0)) + 1}
	// intializing bitset to binary representation of 2^63 + 2^30 + 1, so bits 0, 30, and 63 should be set
	bs := &BitSet{size: 64, words: words}

	bitsToFlip := []int{0, 15, 30, 45, 63}
	bs.FlipBits(bitsToFlip)
	bools, numSet := bs.TestBits(bitsToFlip)

	if numSet != 2 {
		t.Errorf("BitSet.TestBits() want %d, got %d", 2, numSet)
	}
	if !slices.Equal(bools, []bool{false, true, false, true, false}) {
		t.Errorf("BitSet.TestBits() in ClearBits: want [true, false, true, false, true], got %v", bools)
	}
}

func TestBitSet_Or_EqualLength(t *testing.T) {
	a := NewBitSetInitialSize(10)
	b := NewBitSetInitialSize(90)

	aToSet := []int{1}
	bToSet := []int{0, 2, 4, 6, 45}
	a.SetBits(aToSet)
	b.SetBits(bToSet)
	a.Or(b)
	fmt.Println(a)
}

func TestBitSet_Or_SmallerReceiver(t *testing.T) {
	a := NewBitSetInitialSize(80)
	b := NewBitSetInitialSize(200)

	aToSet := []int{1}
	bToSet := []int{0, 2, 4, 6, 64}
	a.SetBits(aToSet)
	b.SetBits(bToSet)
	str := b.String()
	fmt.Println("b:", str, len(str), b.words)
	//fmt.Println("a:", a, len(a.String()), "b:", b, len(b.String()), b.words)
	a.Or(b)
	fmt.Println(a)
}

func TestBitSet_Or_LargerReceiver(t *testing.T) {
	a := NewBitSetInitialSize(10)
	b := NewBitSetInitialSize(20)

	aToSet := []int{1}
	bToSet := []int{0, 2, 4, 6, 18}
	a.SetBits(aToSet)
	b.SetBits(bToSet)
	fmt.Println("a:", a, "b:", b)
	a.Or(b)
	fmt.Println(a)
}

func TestBitSet_And_LargerReceiver(t *testing.T) {
	a := NewBitSetInitialSize(80)
	b := NewBitSetInitialSize(156)

	aToSet := []int{1, 5, 10, 15, 17, 29}
	bToSet := []int{1, 29, 150}
	a.SetBits(aToSet)
	b.SetBits(bToSet)
	str := b.String()
	fmt.Println("b:", str, len(str), b.words)
	a.And(b)
	//res := And(a, b)
	fmt.Println(a)

}

func TestBitSet_And_Smaller(t *testing.T) {
	a := NewBitSetInitialSize(20)
	b := NewBitSetInitialSize(64)

	aToSet := []int{1, 5, 10, 15, 17}
	bToSet := []int{1, 5, 10, 15, 17}
	a.SetBits(aToSet)
	b.SetBits(bToSet)
	fmt.Println("b arr:", b.words)
	fmt.Println("a:", a, "b: ", b)
	a.And(b)
	//res := And(a, b)
	fmt.Printf("%b\n", 0b00101000010000100010&0b1001001001001)
	fmt.Println(a)
}

func TestBitSet_Xor_Smaller(t *testing.T) {
	a := NewBitSetInitialSize(20)
	b := NewBitSetInitialSize(64)

	a.Not()
	b.Not()
	aToClear := []int{1, 5, 10, 15, 17}
	a.ClearBits(aToClear)
	fmt.Println(a, b)
	a.Xor(b)
	fmt.Println(a)
}

func TestBitSet_Not(t *testing.T) {
	a := NewBitSetInitialSize(20)

	aToSet := []int{1, 5, 10, 15, 17}
	a.SetBits(aToSet)
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
	bs := NewBitSetInitialSize(numBits)
	bs.SetBits(bits)
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
	bs := NewBitSetInitialSize(0)
	count, want := bs.CountSetBits(), 0
	if count != want {
		t.Errorf("BitSet.CountSetBits() with empty bitset: got %d, want %d", count, want)
	}
	numBits := 512
	bs, setBits := NewBitSetInitialSize(numBits), make(map[uint64]bool)
	numBitsToSet := rand.Intn(numBits)
	bits := make([]int, numBitsToSet)
	for i := 0; i < numBitsToSet; i++ {
		bit := rand.Intn(numBits)
		setBits[uint64(bit)], bits[i] = true, bit
	}
	bs.SetBits(bits)
	count, want = bs.CountSetBits(), len(setBits)
	if count != want {
		t.Errorf("BitSet has count %d, want %d", count, want)
	}
}

func Test_Do(t *testing.T) {
	fmt.Printf("%b\n", 0b00000|0b1001)
}
