package ops

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/bits"
	"reflect"
)

const (
	MAX_SLICE_SIZE = 1000
)

// Converts a big-endian byte array to uint64
func bytesToUint64(input []byte) (uint64, error) {
	input = trimLeadingZeros(input)
	if !fitsInUint64(input) {
		return 0, fmt.Errorf("%v is not representable as uint64", input)
	}
	num := uint64(0)
	for i, b := range ByteReverse(input) {
		num += uint64(b) << (8 * i)
	}
	return num, nil
}

// Adds two bytes, returning the sum (i.e.: (a+b)%256) and carry (i.e.: (a+b)/256)
func byteAdd(a, b byte) (byte, byte) {
	// Do regular sum modulo 256
	sum := a + b

	// If it overflowed we want a carry
	// In that case the sum will be less than at least one of the inputs
	var carry byte = 0
	if sum < a || sum < b {
		carry = 1
	}

	return sum, carry
}

// Returns the first index such that input[index] != 0, or len(input) if no such input is found
func firstNonZeroIndex(input []byte) uint {
	for i := uint(0); i < Ulen(input); i++ {
		if input[i] != 0 {
			return i
		}
	}
	return Ulen(input)
}

// Returns whether the result can be represented as a uint64
func fitsInUint64(a []byte) bool {
	return nbitsAsUint64(a) <= 64
}

func nbitsAsUint64(input []byte) uint64 {
	trimmed := trimLeadingZeros(input)
	sumUint64 := uint64(0)
	if len(trimmed) > 0 {
		// Since the input has been trimmed, the result must be the number of bits
		// to represent the first byte plus the bit length of everything else
		sumUint64 += uint64(NBitsByte(trimmed[0])) + uint64((len(trimmed)-1)*8)
	}
	return sumUint64
}

// Adds value to the front of the slice
func prepend(slice []byte, value byte) []byte {
	// Allocate an extra slot for the copy operation
	result := append(slice, 0)

	// Move everything one step to the right
	copy(result[1:], result)

	// Then we can fill in the new value
	result[0] = value
	return result
}

// Returns the input with leading zeros removed
func trimLeadingZeros(input []byte) []byte {
	return input[firstNonZeroIndex(input):]
}

// Converts a uint64 to a big-endian byte array
func uint64ToBytes(input uint64) []byte {
	if input == 0 {
		// Canonical zero
		return []byte{}
	} else if input < 256 {
		// Speed hack for small inputs
		return []byte{byte(input)}
	} else {
		// In the general case, defer to binary.BigEndian
		answer := make([]byte, 8)
		binary.BigEndian.PutUint64(answer, input)
		return trimLeadingZeros(answer)
	}
}

// Adds a and b, both representing big-endian numbers
func Add(a, b []byte) []byte {
	a, b = PadToEqualSize(a, b)

	// We will iterate from LSB to MSB for the carry to work
	a = ByteReverse(a)
	b = ByteReverse(b)
	answer := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		// Get the sum and carry of the a and b parts
		sum, carry := byteAdd(a[i], b[i])

		// Add the sum to the existing value in this slot, which
		// may also overflow in which case we need an additional carry
		previousCarryPlusSum, extraCarry := byteAdd(answer[i], sum)
		carry += extraCarry
		answer[i] = previousCarryPlusSum

		// Add the carry to the next slot, or create a slot for it if this is the last one
		if i+1 >= len(answer) {
			if carry > 0 {
				answer = append(answer, carry)
			}
		} else {
			answer[i+1] += carry
		}
	}

	// Flip back to big endian
	return ByteReverse(answer)
}

// Returns bitwise a AND b
func And(a, b []byte) []byte {
	return BinaryOp(a, b, func(ai, bi byte) byte { return ai & bi })
}

// Returns bitwise OP(a, b)
func BinaryOp(a, b []byte, elemfunc func(ai, bi byte) byte) []byte {
	a, b = PadToEqualSize(a, b)
	c := make([]byte, len(a))
	for i, ai := range a {
		c[i] = elemfunc(ai, b[i])
	}
	return c
}

// Returns the binary string representation of the bytes
func BytesToBin(a []byte, nBytes uint) string {
	parts := []byte("0b")

	// Iterate over each byte MSB first
	for _, b := range a {
		// Iterate over each bit in the byte MSb first
		for i := 7; i >= 0; i-- {
			// If the bit is set, emit "1", else emit "0"
			mask := byte(1 << i)
			masked := mask & b
			if masked == 0 {
				parts = append(parts, []byte("0")...)
			} else {
				parts = append(parts, []byte("1")...)
			}
		}
	}
	return string(parts)
}

// Returns the decimal string representation of the bytes
func BytesToDec(a []byte, nBytes uint) string {
	return big.NewInt(0).SetBytes(a).String()
}

// Returns the hexadecimal string representation of the bytes
func BytesToHex(a []byte, nBytes uint) string {
	// Pad with zeros up to nBytes or len(a), whichever is greater
	if nBytes > Ulen(a) {
		a = PrependZeros(a, nBytes-Ulen(a))
	}
	return "0x" + hex.EncodeToString(a)
}

// Returns the input with the byte order reversed
// Leading zeros in the output (due to trailing zeros in the input) are NOT removed
func ByteReverse(input []byte) []byte {
	reversed := make([]byte, len(input))
	copy(reversed, input)
	for i := 0; i < len(reversed)/2; i++ {
		j := len(reversed) - i - 1
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	return reversed
}

// Returns the number of leading zeros in the input
func Clz(input []byte) []byte {
	answerInt := uint64(0)

	for _, b := range input {
		if b == 0 {
			// As long as we see zero-byte just keep adding 8
			answerInt += 8
		} else {
			// Once we encounter a non-zero byte, add up the leading zeroes in that and we are done
			answerInt += uint64(bits.LeadingZeros8(b))
			break
		}
	}
	return uint64ToBytes(answerInt)
}

// Returns whether the two are equivalent except for any leading zeros
func Equivalent(left, right []byte) bool {
	left, right = PadToEqualSize(left, right)
	return reflect.DeepEqual(left, right)
}

// Returns whether the left argument represents a strictly greater number than the right argument
func LeftIsGreater(left, right []byte) bool {
	left, right = PadToEqualSize(left, right)

	// On the first byte that differs we know the result
	for i, l := range left {
		r := right[i]
		if l > r {
			return true
		} else if l < r {
			return false
		}
	}

	// They are equal
	return false
}

// Returns whether the left arguments represents a greater or equal number to the right argument
func LeftIsGreaterOrEqual(left, right []byte) bool {
	left, right = PadToEqualSize(left, right)

	// Zeros are equal
	if len(left) == 0 {
		return true
	}

	// On the first byte that differs we know the result
	for i, l := range left {
		r := right[i]
		if l >= r {
			return true
		} else if l < r {
			return false
		}
	}

	panic("Invalid result for LeftIsGreaterOrEqual")
}

// Returns the number of bits needed to represent the input
func Nbits(input []byte) []byte {
	return uint64ToBytes(nbitsAsUint64(input))
}

// Returns ~a
func Not(input []byte) []byte {
	return UnaryOp(input, func(ai byte) byte { return byte(0xff - int(ai)) })
}

// Returns a OR b
func Or(a, b []byte) []byte {
	return BinaryOp(a, b, func(ai, bi byte) byte { return ai | bi })
}

// Returns slices of equal length representing the same big-endian numbers as a and b
func PadToEqualSize(a, b []byte) ([]byte, []byte) {
	aLen := Ulen(a)
	bLen := Ulen(b)
	if aLen == bLen {
		return a, b
	} else if aLen > bLen {
		return a, PrependZeros(b, aLen-bLen)
	} else {
		return PrependZeros(a, bLen-aLen), b
	}
}

// Returns the popcount of the input
func Popcount(input []byte) []byte {
	answerInt := uint64(0)
	for _, b := range input {
		answerInt += uint64(bits.OnesCount(uint(b)))
	}
	return uint64ToBytes(answerInt)
}

// Prepends n zeros to the slice
// n is capped at MAX_SLICE_SIZE to avoid running out of memory
func PrependZeros(slice []byte, n uint) []byte {
	return append(make([]byte, uintmin(n, MAX_SLICE_SIZE)), slice...)
}

// Returns whether the len(right) last bytes of left are equal to right
func RightIsSuffixOfLeft(left, right []byte) bool {
	// Impossible if right is longer than left
	if len(left) < len(right) {
		return false
	}

	// Check one by one starting from where the right would occur in left
	offset := len(left) - len(right)
	for i, b := range right {
		if b != left[i+offset] {
			return false
		}
	}

	// No differences found
	return true
}

// Returns a >> b
func ShiftLeft(a, b []byte) []byte {
	if len(a) == 0 {
		return a
	}
	b = trimLeadingZeros(b)
	nBits, err := bytesToUint64(b)
	if err != nil || nBits >= nbitsAsUint64(a) {
		return Zeros(Ulen(a))
	}

	// First do the part of the shift that is divisible by 8
	nBytes := nBits / 8

	// Copy over starting from the end of the slice
	for i := uint64(0); i < U64len(a)-nBytes; i++ {
		dest := U64len(a) - 1 - i
		src := dest - nBytes
		a[dest] = a[src]
	}

	// Fill the front with zeros
	for i := uint64(0); i < nBytes; i++ {
		a[i] = 0
	}

	// Then do the remainder
	nBits = nBits % 8
	for i := uint64(0); i < nBits; i++ {
		carry := false
		for j, ai := range a {
			a[j] >>= 1
			if carry {
				a[j] += 1 << 7
			}
			carry = ai&1 == 1
		}
	}

	return a
}

// Parses the input string to a byte array
func StringToBytes(a string) ([]byte, bool) {
	resultInt, ok := big.NewInt(0).SetString(a, 0)
	result := []byte{}
	if ok {
		result = resultInt.Bytes()
	}
	return result, ok
}

// Subtracts b from a, both representing big-endian numbers
func Subtract(a, b []byte) []byte {
	a, b = PadToEqualSize(a, b)
	subtracted := Add(a, TwosComplement(b))
	return Truncate(subtracted, Ulen(a))
}

// Returns the last n bytes in a
func Truncate(a []byte, n uint) []byte {
	if n >= Ulen(a) {
		return a
	} else {
		return a[Ulen(a)-n:]
	}
}

// Returns -a
func TwosComplement(input []byte) []byte {
	// Special case for 0-bit zero
	if len(input) == 0 {
		return input
	}

	// Otherwise -a is ~a+1
	return Truncate(Add(Not(input), []byte{1}), Ulen(input))
}

// Returns the unsigned length of the input
func Ulen(a []byte) uint {
	return uint(len(a))
}

// Returns OP(a)
func UnaryOp(a []byte, elemfunc func(ai byte) byte) []byte {
	c := make([]byte, len(a))
	for i, ai := range a {
		c[i] = elemfunc(ai)
	}
	return c
}

// Returns the uint64 length of the input
func U64len(a []byte) uint64 {
	return uint64(len(a))
}

// Returns a XOR b
func Xor(a, b []byte) []byte {
	return BinaryOp(a, b, func(ai, bi byte) byte { return ai ^ bi })
}

// Returns zero-initialized byte slice of length n
func Zeros(n uint) []byte {
	return make([]byte, n)
}
