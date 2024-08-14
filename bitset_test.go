package bitset

import (
	"fmt"
	"math"
	"testing"
)

func TestBitset_BitArrayLen(t *testing.T) {
	bs := NewBitset(4416)
	want := 69
	if len(bs.bits) != 69 {
		t.Errorf("Bitset has length %d, want %d", len(bs.bits), want)
	}

	bs = NewBitset(420)
	want = 7
	if len(bs.bits) != 7 {
		t.Errorf("Bitset has length %d, want %d", len(bs.bits), want)
	}
}

func TestBitset_OutOfBounds(t *testing.T) {
	bs := NewBitset(4416)
	err := bs.Set(4416)
	if err == nil {
		t.Errorf("Bitset should have returned an error")
	}

	err = bs.Set(-1)
	if err == nil {
		t.Errorf("Bitset should have returned an error")
	}
}

func TestBitset_Set(t *testing.T) {
	numBits := 512
	bs := NewBitset(numBits)
	for i := len(bs.bits) - 1; i >= 0; i-- {
		bitPos := i + (i * 64)
		err := bs.Set(bitPos)
		if err != nil {
			t.Error(err)
		}
		idx := len(bs.bits) - 1 - bitPos/64
		if bs.bits[idx] != uint64(math.Pow(2.0, float64(i))) {
			t.Errorf("Set failed for word %d at index %d", bs.bits[i], i)
		}
	}
}

func TestBitset_Test(t *testing.T) {

}

func TestBitset_Clear(t *testing.T) {}

func TestBitset_Flip(t *testing.T) {}

func TestBitset_Not(t *testing.T) {}

func TestBitset_String(t *testing.T) {
	bs := NewBitset(10)
	var err error
	err = bs.Set(5)
	err = bs.Set(3)
	err = bs.Set(1)
	//fmt.Println(bs.String())
	fmt.Println(err)
}

func TestBitset_Count(t *testing.T) {
	bs := NewBitset(64)
	want := 100
	var err error
	for i := 0; i < want; i++ {
		err = bs.Set(i)
	}
	//fmt.Println(bs.String())
	if err != nil {
		t.Error(err)
	}
	count := bs.Count()
	if count != want {
		t.Errorf("Bitset has count %d, want %d", count, want)
	}
}
