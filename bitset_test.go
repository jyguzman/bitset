package bitset

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestBitSet_BitArrayLen(t *testing.T) {
	bs := NewBitSet(4416)
	want := 69
	if len(bs.bitArray) != 69 {
		t.Errorf("BitSet has length %d, want %d", len(bs.bitArray), want)
	}

	bs = NewBitSet(420)
	want = 7
	if len(bs.bitArray) != 7 {
		t.Errorf("BitSet has length %d, want %d", len(bs.bitArray), want)
	}
}

func TestBitSet_OutOfBounds(t *testing.T) {
	bs := NewBitSet(4416)
	err := bs.Set(4416)
	if err == nil {
		t.Errorf("BitSet should have returned an error")
	}

	err = bs.Set(-1)
	if err == nil {
		t.Errorf("BitSet should have returned an error")
	}
}

func TestBitSet_Set(t *testing.T) {
	numBits := 544
	bs := NewBitSet(numBits)
	for i := len(bs.bitArray) - 1; i >= 0; i-- {
		bitPos := i + (i * 64)
		err := bs.Set(bitPos)
		if err != nil {
			t.Error(err)
		}
		idx := len(bs.bitArray) - 1 - bitPos/64
		if bs.bitArray[idx] != uint64(math.Pow(2.0, float64(i))) {
			t.Errorf("Set failed for word %d at index %d", bs.bitArray[i], i)
		}
	}
}

func TestBitSet_Test(t *testing.T) {
	numBits := 544
	bs := NewBitSet(numBits)
	for i := len(bs.bitArray) - 1; i >= 0; i-- {
		bitPos := i + (i * 64)
		if err := bs.Set(bitPos); err != nil {
			t.Error(err)
		}
		idx := len(bs.bitArray) - 1 - bitPos/64
		isSet, err := bs.Test(idx)
		if err != nil {
			t.Error(err)
		}
		if !isSet {
			t.Errorf("BitSet test failed for bit %d", bitPos)
		}
	}
}

func TestBitSet_Clear(t *testing.T) {}

func TestBitSet_Flip(t *testing.T) {}

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
	b := NewBitSet(20)

	aToSet := []int{0, 1, 2, 4}
	bToSet := []int{0, 1, 2, 4}
	if err := a.SetBits(aToSet); err != nil {
		t.Error(err)
	}
	if err := b.SetBits(bToSet); err != nil {
		t.Error(err)
	}
	res := a.And(b)
	fmt.Println(res)
}

func TestBitSet_String(t *testing.T) {
	bs := NewBitSet(512)
	var err error
	err = bs.Set(132)
	err = bs.Set(65)
	err = bs.Set(1)
	fmt.Println(bs.bitArray)
	fmt.Println(bs.String())
	fmt.Println(err)
}

func TestBitSet_Count(t *testing.T) {
	numBits := 512
	bs := NewBitSet(numBits)
	setBits := make(map[uint64]bool)
	var bits []int
	numBitsToSet := rand.Intn(numBits)
	for i := 0; i < numBitsToSet; i++ {
		bit := rand.Intn(numBits)
		setBits[uint64(bit)] = true
		bits = append(bits, bit)
	}
	if err := bs.SetBits(bits); err != nil {
		t.Error(err)
	}
	count := bs.CountSetBits()
	want := len(setBits)
	if count != want {
		t.Errorf("BitSet has count %d, want %d", count, want)
	}
	fmt.Println(bs.String())
}
