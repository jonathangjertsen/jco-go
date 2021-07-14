package ops

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"math/bits"
	"reflect"
)

const (
	MAX_SLICE_SIZE = 1000
)

// Adds two bytes, returning the sum (i.e.: (a+b)%256) and carry (i.e.: (a+b)/256)
func byteAdd(a, b byte) (byte, byte) {
	sum := a + b
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

// Returns whichever is greatest of a and b
func intmax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// Adds value to the front of the slice
func prepend(slice []byte, value byte) []byte {
	result := append(slice, 0)
	copy(result[1:], result)
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
		return []byte{}
	} else if input < 256 {
		return []byte{byte(input)}
	} else {
		answer := make([]byte, 8)
		binary.BigEndian.PutUint64(answer, input)
		return trimLeadingZeros(answer)
	}
}

// Adds a and b, both representing big-endian numbers
func Add(a, b []byte) []byte {
	a, b = PadToEqualSize(a, b)
	a = ByteReverse(a)
	b = ByteReverse(b)
	answer := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		sum, carry := byteAdd(a[i], b[i])
		previousCarryPlusSum, extraCarry := byteAdd(answer[i], sum)
		carry += extraCarry
		answer[i] = previousCarryPlusSum
		if i+1 >= len(answer) {
			if carry > 0 {
				answer = append(answer, carry)
			}
		} else {
			answer[i+1] += carry
		}
	}
	return ByteReverse(answer)
}

// Returns a AND b
func And(a, b []byte) []byte {
	return BinaryOp(a, b, func(ai, bi byte) byte { return ai & bi })
}

// Returns OP(a, b)
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
	for _, b := range a {
		for i := 7; i >= 0; i-- {
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
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	return reversed
}

// Returns the number of leading zeros in the input
func Clz(input []byte) []byte {
	answerInt := uint64(0)
	for _, b := range input {
		if b == 0 {
			answerInt += 8
		} else {
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

func LeftIsGreater(left, right []byte) bool {
	left, right = PadToEqualSize(left, right)
	for i, l := range left {
		r := right[i]
		if l > r {
			return true
		} else if l < r {
			return false
		}
	}
	return false
}

func LeftIsGreaterOrEqual(left, right []byte) bool {
	left, right = PadToEqualSize(left, right)
	if len(left) == 0 {
		return true
	}
	for i, l := range left {
		r := right[i]
		if l >= r {
			return true
		} else if l < r {
			return false
		}
	}
	return false
}

// Returns the number of bits needed to represent the input
func Nbits(input []byte) []byte {
	trimmed := trimLeadingZeros(input)
	sumUint64 := uint64(0)
	if len(trimmed) > 0 {
		sumUint64 += uint64(NBitsByte(trimmed[0])) + uint64((len(trimmed)-1)*8)
	}
	return uint64ToBytes(sumUint64)
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
func PrependZeros(slice []byte, n uint) []byte {
	if n > MAX_SLICE_SIZE {
		n = MAX_SLICE_SIZE
	}

	zeros := make([]byte, n)
	return append(zeros, slice...)
}

func RightIsSuffixOfLeft(left, right []byte) bool {
	if len(left) < len(right) {
		return false
	}
	offset := len(left) - len(right)
	for i, b := range right {
		if b != left[i+offset] {
			return false
		}
	}
	return true
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
	if len(input) == 0 {
		return input
	}
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
