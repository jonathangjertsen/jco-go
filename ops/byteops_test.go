package ops

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
	"testing/quick"
)

func check(t *testing.T, f interface{}) {
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("%v\n\nStack:%s\n", err, debug.Stack())
	}
}

func leftIsGreater(left, right []byte) bool {
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

func leftIsGreaterOrEqual(left, right []byte) bool {
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

func rightIsSuffixOfLeft(left, right []byte) bool {
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

func TestByteAdd(t *testing.T) {
	// Vector part
	var vector = []struct {
		a         byte
		b         byte
		wantSum   byte
		wantCarry byte
	}{
		{0, 0, 0, 0},
		{0xf0, 0x0f, 0xff, 0x00},
		{0xf0, 0x10, 0x00, 0x01},
		{0xff, 0xff, 0xfe, 0x01},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v,%v\n", tt.a, tt.b)
		t.Run(testname, func(t *testing.T) {
			haveSum, haveCarry := byteAdd(tt.a, tt.b)
			if haveSum != tt.wantSum || haveCarry != tt.wantCarry {
				t.Errorf(
					"Sum: want sum=%v, carry=%v, have sum=%v, carry=%v\n",
					tt.wantSum,
					tt.wantCarry,
					haveSum,
					haveCarry,
				)
			}
		})
	}

	// Property: carry is 1 or 0
	check(t, func(a, b byte) bool {
		_, carry := byteAdd(a, b)
		return (carry == 1) || (carry == 0)
	})

	// Property: byteAdd is associative
	check(t, func(a, b byte) bool {
		sum1, carry1 := byteAdd(a, b)
		sum2, carry2 := byteAdd(b, a)
		return sum1 == sum2 && carry1 == carry2
	})

	// Property: byteAdd with 0 returns sum=a and carry=0
	check(t, func(a byte) bool {
		sum, carry := byteAdd(a, 0)
		return sum == a && carry == 0
	})
}

func TestByteReverse(t *testing.T) {
	var vector = []struct {
		input []byte
		want  []byte
	}{
		{
			[]byte{0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
		},
		{
			[]byte{0xff, 0},
			[]byte{0, 0xff},
		},
		{
			[]byte{0, 0xff},
			[]byte{0xff, 0},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v\n", tt.input)
		t.Run(testname, func(t *testing.T) {
			have := ByteReverse(tt.input)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: reverse twice yields the input
	check(t, func(input []byte) bool {
		doublyReversed := ByteReverse(ByteReverse(input))
		ok := reflect.DeepEqual(doublyReversed, input)
		if !ok {
			fmt.Printf("Failing doublyReversed: %v\n", doublyReversed)
		}
		return ok
	})

	// Property: if input is not empty, the first in the reversed equals the last in the input
	check(t, func(input []byte) bool {
		reversed := ByteReverse(input)
		ok := len(input) == 0 || input[0] == reversed[len(input)-1]
		if !ok {
			fmt.Printf("Failing reversed: %v\n", reversed)
		}
		return ok
	})

	// Property: reverse contains the same bytes as the input
	check(t, func(input []byte) bool {
		inputReversed := ByteReverse(input)
		inputByteCount := map[byte]int{}
		inputReversedByteCount := map[byte]int{}
		for i, b := range input {
			inputByteCount[b]++
			inputReversedByteCount[inputReversed[i]]++
		}
		ok := reflect.DeepEqual(inputByteCount, inputReversedByteCount)
		if !ok {
			fmt.Printf(
				"Failing inputByteCount: %v, inputReversedByteCount: %v\n",
				inputByteCount,
				inputReversedByteCount,
			)
		}
		return ok
	})
}

func TestTrimLeadingZeros(t *testing.T) {
	var vector = []struct {
		input []byte
		want  []byte
	}{
		{
			[]byte{0, 0, 0, 0},
			[]byte{},
		},
		{
			[]byte{0xff, 0},
			[]byte{0xff, 0},
		},
		{
			[]byte{0, 0xff},
			[]byte{0xff},
		},
		{
			[]byte{0xff, 0xff},
			[]byte{0xff, 0xff},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v\n", tt.input)
		t.Run(testname, func(t *testing.T) {
			have := trimLeadingZeros(tt.input)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: trimmed will be a smaller slice or exactly equal to input
	check(t, func(a []byte) bool {
		trimmed := trimLeadingZeros(a)
		ok := len(trimmed) < len(a) || reflect.DeepEqual(trimmed, a)
		if !ok {
			fmt.Printf("Failing trimmed: %v\n", trimmed)
		}
		return ok
	})

	// Property: trimmed is a suffix of the original
	check(t, func(a []byte) bool {
		trimmed := trimLeadingZeros(a)
		ok := rightIsSuffixOfLeft(a, trimmed)
		if !ok {
			fmt.Printf("Failing trimmed: %v\n", trimmed)
		}
		return ok
	})

	// Property: trimmed is equivalent to input
	check(t, func(a []byte) bool {
		trimmed := trimLeadingZeros(a)
		ok := Equivalent(a, trimmed)
		if !ok {
			fmt.Printf("Failing trimmed: %v\n", trimmed)
		}
		return ok
	})
}

func TestAdd(t *testing.T) {
	var vector = []struct {
		a    []byte
		b    []byte
		want []byte
	}{
		{
			[]byte{0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
		},
		{
			[]byte{0x10, 0x10},
			[]byte{0x20, 0x20},
			[]byte{0x30, 0x30},
		},
		{
			[]byte{0xff, 0xff},
			[]byte{0xff, 0xff},
			[]byte{0x01, 0xff, 0xfe},
		},
		{
			[]byte{0xff, 0xff},
			[]byte{0, 0, 0, 0, 0, 0, 0xff, 0xff},
			[]byte{0, 0, 0, 0, 0, 0x01, 0xff, 0xfe},
		},
		{
			[]byte{20, 110, 177, 39, 171, 32, 91, 192, 152},
			[]byte{235, 201, 191, 208, 116, 20, 104, 155, 109},
			[]byte{1, 0, 56, 112, 248, 31, 52, 196, 92, 5},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v,%v\n", tt.a, tt.b)
		t.Run(testname, func(t *testing.T) {
			have := Add(tt.a, tt.b)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: len(a+b) >= len(a)
	check(t, func(a, b []byte) bool {
		added := Add(a, b)
		ok := len(added) >= len(a) && len(added) >= len(b)
		if !ok {
			fmt.Printf("Failing added: %v\n", added)
		}
		return ok
	})

	// Property: a+b >= a and a+b >= b
	check(t, func(a, b []byte) bool {
		added := Add(a, b)
		ok := leftIsGreaterOrEqual(added, a) && leftIsGreaterOrEqual(added, b)
		if !ok {
			fmt.Printf("Failing added: %v\n", added)
		}
		return ok
	})

	// Property: a + b == b + a
	check(t, func(a, b []byte) bool {
		apb := Add(a, b)
		bpa := Add(b, a)
		ok := reflect.DeepEqual(apb, bpa)
		if !ok {
			fmt.Printf("Failing a+b: %v, b+a: %v\n", apb, bpa)
		}
		return ok
	})

	// Property: a + 0 == a
	check(t, func(a []byte) bool {
		ok := reflect.DeepEqual(a, Add(a, []byte{})) && reflect.DeepEqual(a, Add([]byte{}, a)) && Equivalent(a, Add([]byte{0, 0, 0}, a))
		if !ok {
			fmt.Printf("Failing value: %v\n", a)
		}
		return ok
	})
}

func TestSubtract(t *testing.T) {
	var vector = []struct {
		a    []byte
		b    []byte
		want []byte
	}{
		{
			[]byte{},
			[]byte{},
			[]byte{},
		},
		{
			[]byte{0},
			[]byte{0},
			[]byte{0},
		},
		{
			[]byte{0x20, 0x20},
			[]byte{0x10, 0x10},
			[]byte{0x10, 0x10},
		},
		{
			[]byte{0xff, 0xff},
			[]byte{0xff, 0xff},
			[]byte{0x00, 0x00},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v,%v\n", tt.a, tt.b)
		t.Run(testname, func(t *testing.T) {
			have := Subtract(tt.a, tt.b)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: a - b == -(b - a)
	check(t, func(a, b []byte) bool {
		amb := Subtract(a, b)
		bma := Subtract(b, a)
		bmaTwosComplement := TwosComplement(bma)
		ok := reflect.DeepEqual(amb, bmaTwosComplement)
		if !ok {
			fmt.Printf("Failing a-b: %v, -(b-a): %v\n", amb, bmaTwosComplement)
		}
		return ok
	})

	// Property: a - 0 == a
	check(t, func(a []byte) bool {
		subtractingEmptyZero := Subtract(a, []byte{})
		subtractingNonEmptyZero := Subtract(a, []byte{0, 0, 0})
		ok := reflect.DeepEqual(a, subtractingEmptyZero) && Equivalent(a, subtractingNonEmptyZero)
		if !ok {
			fmt.Printf("Failing subtractingEmptyZero: %v, subtractingNonEmptyZero: %v\n", subtractingEmptyZero, subtractingNonEmptyZero)
		}
		return ok
	})

	// Property: 0 - a == -a
	check(t, func(a []byte) bool {
		subtractingFromEmptyZero := Subtract([]byte{}, a)
		ok := reflect.DeepEqual(TwosComplement(a), subtractingFromEmptyZero)
		if !ok {
			fmt.Printf("Failing subtractingFromEmptyZero: %v\n", subtractingFromEmptyZero)
		}
		return ok
	})
}

func TestOr(t *testing.T) {
	var vector = []struct {
		a    []byte
		b    []byte
		want []byte
	}{
		{
			[]byte{0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
		},
		{
			[]byte{0x10, 0x10},
			[]byte{0x20, 0x20},
			[]byte{0x30, 0x30},
		},
		{
			[]byte{0xff, 0xff},
			[]byte{0xff, 0xff},
			[]byte{0xff, 0xff},
		},
		{
			[]byte{0xff, 0xff},
			[]byte{0, 0, 0, 0, 0, 0, 0xff, 0xff},
			[]byte{0, 0, 0, 0, 0, 0, 0xff, 0xff},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v,%v\n", tt.a, tt.b)
		t.Run(testname, func(t *testing.T) {
			have := Or(tt.a, tt.b)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: len(a|b) == len(a)
	check(t, func(a, b []byte) bool {
		ored := Or(a, b)
		a, b = PadToEqualSize(a, b)
		return len(ored) == len(a) && len(ored) == len(b)
	})

	// Property: a | 0 == a
	check(t, func(a []byte) bool {
		zero := make([]byte, len(a))
		for i, _ := range zero {
			zero[i] = 0
		}
		return Equivalent(a, Or(a, zero))
	})

	// Property: a | -1 == -1
	check(t, func(a []byte) bool {
		minus1 := make([]byte, len(a))
		for i, _ := range minus1 {
			minus1[i] = 0xff
		}
		return Equivalent(minus1, Or(a, minus1))
	})
}

func TestAnd(t *testing.T) {
	var vector = []struct {
		a    []byte
		b    []byte
		want []byte
	}{
		{
			[]byte{0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
		},
		{
			[]byte{0x10, 0x10},
			[]byte{0x20, 0x20},
			[]byte{0x00, 0x00},
		},
		{
			[]byte{0xff, 0xff},
			[]byte{0xff, 0xff},
			[]byte{0xff, 0xff},
		},
		{
			[]byte{0xff, 0xff},
			[]byte{0, 0, 0, 0, 0, 0, 0xff, 0xff},
			[]byte{0, 0, 0, 0, 0, 0, 0xff, 0xff},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v,%v\n", tt.a, tt.b)
		t.Run(testname, func(t *testing.T) {
			have := And(tt.a, tt.b)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: len(a&b) == len(a)
	check(t, func(a, b []byte) bool {
		ored := And(a, b)
		a, b = PadToEqualSize(a, b)
		return len(ored) == len(a) && len(ored) == len(b)
	})

	// Property: a & 0 == 0
	check(t, func(a []byte) bool {
		zero := make([]byte, len(a))
		for i, _ := range zero {
			zero[i] = 0
		}
		return Equivalent(zero, And(a, zero))
	})

	// Property: a & -1 == a
	check(t, func(a []byte) bool {
		minus1 := make([]byte, len(a))
		for i, _ := range minus1 {
			minus1[i] = 0xff
		}
		return Equivalent(a, And(a, minus1))
	})
}

func TestXor(t *testing.T) {
	var vector = []struct {
		a    []byte
		b    []byte
		want []byte
	}{
		{
			[]byte{0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
		},
		{
			[]byte{0x10, 0x10},
			[]byte{0x20, 0x20},
			[]byte{0x30, 0x30},
		},
		{
			[]byte{0xff, 0xff},
			[]byte{0xff, 0xff},
			[]byte{0x00, 0x00},
		},
		{
			[]byte{0xff, 0xff},
			[]byte{0, 0, 0, 0, 0, 0, 0xff, 0xff},
			[]byte{0, 0, 0, 0, 0, 0, 0x00, 0x00},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v,%v\n", tt.a, tt.b)
		t.Run(testname, func(t *testing.T) {
			have := Xor(tt.a, tt.b)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: len(a^b) == len(a)
	check(t, func(a, b []byte) bool {
		xored := Xor(a, b)
		a, b = PadToEqualSize(a, b)
		return len(xored) == len(a) && len(xored) == len(b)
	})

	// Property: a ^ 0 == a
	check(t, func(a []byte) bool {
		zero := make([]byte, len(a))
		for i, _ := range zero {
			zero[i] = 0
		}
		return Equivalent(a, Xor(a, zero))
	})

	// Property: a ^ -1 == ~a
	check(t, func(a []byte) bool {
		minus1 := make([]byte, len(a))
		for i, _ := range minus1 {
			minus1[i] = 0xff
		}
		return Equivalent(Not(a), Xor(a, minus1))
	})
}

func TestPopcount(t *testing.T) {
	var vector = []struct {
		input []byte
		want  []byte
	}{
		{
			[]byte{0, 0, 0, 0},
			[]byte{},
		},
		{
			[]byte{0xff, 0},
			[]byte{8},
		},
		{
			[]byte{0xff, 0xaa},
			[]byte{12},
		},
		{
			[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			[]byte{1, 0},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v\n", tt.input)
		t.Run(testname, func(t *testing.T) {
			have := Popcount(tt.input)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: len(popcount(a)) <= max(1, len(a)/32)
	check(t, func(a []byte) bool {
		lpc := len(Popcount(a))
		return lpc == 1 || lpc <= len(a)/32
	})
}

func TestClz(t *testing.T) {
	var vector = []struct {
		input []byte
		want  []byte
	}{
		{
			[]byte{},
			[]byte{},
		},
		{
			[]byte{0x7f, 0},
			[]byte{1},
		},
		{
			[]byte{0xff, 0xaa},
			[]byte{},
		},
		{
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]byte{1, 7},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v\n", tt.input)
		t.Run(testname, func(t *testing.T) {
			have := Clz(tt.input)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}
}

func TestNbits(t *testing.T) {
	var vector = []struct {
		input []byte
		want  []byte
	}{
		{
			[]byte{},
			[]byte{},
		},
		{
			[]byte{0x7f, 0},
			[]byte{15},
		},
		{
			[]byte{0xff, 0xaa},
			[]byte{16},
		},
		{
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]byte{1},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v\n", tt.input)
		t.Run(testname, func(t *testing.T) {
			have := Nbits(tt.input)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}
}

func TestNot(t *testing.T) {
	var vector = []struct {
		input []byte
		want  []byte
	}{
		{
			[]byte{},
			[]byte{},
		},
		{
			[]byte{0xff},
			[]byte{0x00},
		},
		{
			[]byte{0x7f, 0},
			[]byte{0x80, 0xff},
		},
		{
			[]byte{0xff, 0xaa},
			[]byte{0x00, 0x55},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v\n", tt.input)
		t.Run(testname, func(t *testing.T) {
			have := Not(tt.input)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: not twice yields the input
	check(t, func(input []byte) bool {
		doublyInverted := Not(Not(input))
		ok := reflect.DeepEqual(doublyInverted, input)
		if !ok {
			fmt.Printf("Failing doublyInverted: %v\n", doublyInverted)
		}
		return ok
	})
}

func TestTwosComplement(t *testing.T) {
	var vector = []struct {
		input []byte
		want  []byte
	}{
		{
			[]byte{},
			[]byte{},
		},
		{
			[]byte{0xff},
			[]byte{0x01},
		},
		{
			[]byte{0x7f, 0},
			[]byte{0x81, 0x00},
		},
		{
			[]byte{0xff, 0xaa},
			[]byte{0x00, 0x56},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v\n", tt.input)
		t.Run(testname, func(t *testing.T) {
			have := TwosComplement(tt.input)
			if !bytes.Equal(have, tt.want) {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: TwosComplement twice yields the input
	check(t, func(input []byte) bool {
		doublyInverted := TwosComplement(TwosComplement(input))
		ok := reflect.DeepEqual(doublyInverted, input)
		if !ok {
			fmt.Printf("Failing doublyInverted: %v\n", doublyInverted)
		}
		return ok
	})
}

func TestStringToBytes(t *testing.T) {
	var vector = []struct {
		input string
		want  []byte
	}{
		{
			"0x00",
			[]byte{},
		},
		{
			"0xff",
			[]byte{0xff},
		},
		{
			"255",
			[]byte{255},
		},
		{
			"0x222",
			[]byte{0x02, 0x22},
		},
		{
			"0x7f00",
			[]byte{0x7f, 0x00},
		},
		{
			"0xffaa",
			[]byte{0xff, 0xaa},
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v\n", tt.input)
		t.Run(testname, func(t *testing.T) {
			have, ok := StringToBytes(tt.input)
			if !bytes.Equal(have, tt.want) || !ok {
				t.Errorf("Want %v, have %v (ok: %v)", tt.want, have, ok)
			}
		})
	}

	// Property: StringsToBytes does not crash
	check(t, func(a string) bool {
		StringToBytes(a)
		return true
	})
}

func TestBytesToHex(t *testing.T) {
	var vector = []struct {
		input  []byte
		nBytes uint
		want   string
	}{
		{
			[]byte{},
			1,
			"0x00",
		},
		{
			[]byte{0},
			1,
			"0x00",
		},
		{
			[]byte{0, 0},
			1,
			"0x0000",
		},
		{
			[]byte{0, 0},
			2,
			"0x0000",
		},
		{
			[]byte{0, 0},
			3,
			"0x000000",
		},
		{
			[]byte{0xff},
			1,
			"0xff",
		},
		{
			[]byte{0x02, 0x55},
			2,
			"0x0255",
		},
		{
			[]byte{0x02, 0x22},
			2,
			"0x0222",
		},
		{
			[]byte{0x7f, 0x00},
			2,
			"0x7f00",
		},
		{
			[]byte{0x7f, 0x00},
			1,
			"0x7f00",
		},
		{
			[]byte{0x7f, 0x00},
			6,
			"0x000000007f00",
		},
		{
			[]byte{0xff, 0xaa},
			2,
			"0xffaa",
		},
	}
	for _, tt := range vector {
		testname := fmt.Sprintf("%v\n", tt.input)
		t.Run(testname, func(t *testing.T) {
			have := BytesToHex(tt.input, tt.nBytes)
			if have != tt.want {
				t.Errorf("Want %v, have %v\n", tt.want, have)
			}
		})
	}

	// Property: BytesToHex does not truncate
	check(t, func(a []byte, n uint) bool {
		if uint(len(a)) >= n {
			return reflect.DeepEqual(BytesToHex(a, n), BytesToHex(a, uint(len(a))))
		} else {
			return strings.HasSuffix(BytesToHex(a, n)[2:], BytesToHex(a, uint(len(a)))[2:])
		}
	})
}
