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
	bs := NewBitset(512)
	for i := range bs.bits {
		err := bs.Set(i + (i * 64))
		if err != nil {
			panic(err)
		}
		if bs.bits[i] != uint64(math.Pow(2.0, float64(i))) {
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
	fmt.Println(bs.String())
	fmt.Println(err)
}
