package serpent

import (
	"fmt"
	"testing"
)

var bs Bitstring

func init() {
	bs = "11001110010100110010100101010010101110001001100101" +
		"00101001010111001010010101011001001001010110101001010010" +
		"1010010101001010100101"
}

// Tests

// Function TestQuadSplit tests splitting a 128-bit Bitstring
// into a 4 element []Bitstring each of 32-bits.
func TestQuadSplit(t *testing.T) {
	s := bs.QuadSplit()
	sl := len(s)
	if sl == 4 {
		for i := 0; i < sl; i++ {
			sls := len(s[i])
			if sls != 32 {
				t.Fail()
			}
		}
	}
	fmt.Printf("PASS\t-\t\tTestQuadSplit\n")
}

// Function TestQuadJoin tests joining a 4 element []Bitstring of 32-bits
// each into a single 128-bit Bitstring.
func TestQuadJoin(t *testing.T) {
	var ts Bitstring
	s := bs.QuadSplit()
	ts = ts.QuadJoin(s)
	lts := len(ts)
	if lts != 128 {
		t.Fail()
	}
	fmt.Printf("PASS\t-\t\tTestQuadJoin\n")
}

// Function TestLinearTranslation takes a 128-bit Bitstring and performs
// a translation then an inverse translation which should return the
// original Bitstring
func TestLinearTranslation(t *testing.T) {
	LtBitstring := LT(bs)
	LtiBitstring := LTInverse(LtBitstring)
	if LtiBitstring != bs {
		t.Fail()
	}
	fmt.Printf("PASS\t-\t\tTestLinearTranslation\n")
}

// Function TestSHat applies the parallel array of 32 copies of S-Box #3 and
// then applies the inverse to return the Bitstring to it's original form.
func TestSHat(t *testing.T) {
	bs1 := SHat(3, bs)
	orig := SHatInverse(3, bs1)
	if orig != bs {
		t.Fail()
	}
	fmt.Printf("PASS\t-\t\tTestSHat\n")
}

// Function TestLTBitslice applies the equations-based linear transformation
// and then reverses it.
func TestLTBitslice(t *testing.T) {
	bss1 := bs.QuadSplit()
	bsl1 := LTBitslice(bss1)
	bss2 := LTBitsliceInverse(bsl1)
	for i := 0; i < 4; i++ {
		if bss2[i] != bss1[i] {
			t.Fail()
		}
	}
	fmt.Printf("Pass\t-\t\tTestLTBitslice\n")
}

// Examples

// This example takes a bitstring of "000111" and shifts it left by 2 places
// showing the visual effect of the little-endian representation.
func ExampleBitstring_ShiftLeft() {
	var bs1 Bitstring = "000111"
	output := bs1.ShiftLeft(2)
	fmt.Printf("%v\n", output)
	// Output:
	// 000001
}

// This example takes a bitstring of "000111" and shifts it right by 2 places
// showing the visual effect of the little-endian representation.
func ExampleBitstring_ShiftRight() {
	var bs1 Bitstring = "000111"
	output := bs1.ShiftRight(2)
	fmt.Printf("%v\n", output)
	// Output:
	// 011100
}

// This example takes an integer '10' and converts it to a bitstring of at
// least 8 bits long.
func ExampleBitstring_FromInt() {
	var num int = 10
	var bs Bitstring
	output := bs.FromInt(num, 8)
	fmt.Printf("%v\n", output)
	// Output:
	// 01010000
}

// This example demonstrates a Xor operation of two 5-bit bitstrings.
func ExampleBitstring_BinaryXor() {
	var bs1 Bitstring = "10010"
	var bs2 Bitstring = "00011"
	fmt.Printf("%v\n", bs1.BinaryXor(bs2))
	// Output:
	// 10001
}

// This example shows a Xor operation on a 3 element []Bitstring.
func ExampleBitstring_Xor() {
	var bs []Bitstring = []Bitstring{"01", "11", "10"}
	var output Bitstring
	fmt.Printf("%v\n", output.Xor(bs))
	// Output:
	// 00
}

// This example shows a left rotation of 2 places on an 6-bit Bitstring.
func ExampleBitstring_RotateLeft() {
	var bs Bitstring = "000111"
	fmt.Printf("%v\n", bs.RotateLeft(2))
	// Output:
	// 110001
}

// This example shows a right rotation of 2 places on a 6-bit Bitstring.
func ExampleBitstring_RotateRight() {
	var bs Bitstring = "000111"
	fmt.Printf("%v\n", bs.RotateRight(2))
	// Output:
	// 011100
}
