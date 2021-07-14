package ops

func BitsToBytes(bits uint64) uint64 {
	return (bits + 7) / 8
}

// Returns the number of bits needed to represent the input
func NBitsUint64(input uint64) uint64 {
	for i := 63; i >= 0; i-- {
		mask := uint64(1 << i)
		masked := mask & input
		if masked != 0 {
			return uint64(i + 1)
		}
	}
	return 0
}

// Returns the number of bytes needed to represent the input
func NBytesUint64(input uint64) uint64 {
	return BitsToBytes(NBitsUint64(input))
}