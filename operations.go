package bitset

// set sets the Nth bit to 1.
func set(n int, words []uint64) {
	wordIdx, bitIdx := getWordAndPos(n)
	words[wordIdx] |= 1 << bitIdx
}

// reset zeroes the Nth bit.
func reset(n int, words []uint64) {
	wordIdx, bitIdx := getWordAndPos(n)
	words[wordIdx] &= ^(1 << bitIdx)
}

// flip flips the Nth bit, i.e. 0 -> 1 or 1 -> 0.
func flip(n int, words []uint64) {
	wordIdx, bitIdx := getWordAndPos(n)
	words[wordIdx] ^= 1 << bitIdx
}

// test checks if the Nth bit is set.
func test(n int, words []uint64) bool {
	wordIdx, bitIdx := getWordAndPos(n)
	return words[wordIdx]&(1<<bitIdx) >= 1
}

// not flips each bit of the bit array.
func not(size int, words []uint64) {
	bitsLeft := size
	for i := range words {
		words[i] = mask(^words[i], bitsLeft%64)
		bitsLeft -= 64
	}
}

func getWordAndPos(n int) (int, int) {
	return n / 64, n % 64
}
