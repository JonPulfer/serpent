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
	fmt.Printf("Running\t-\t\tTestQuadSplit\n")
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
}

// Function TestQuadJoin tests joining a 4 element []Bitstring of 32-bits
// each into a single 128-bit Bitstring.
func TestQuadJoin(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestQuadJoin\n")
	var ts Bitstring
	s := bs.QuadSplit()
	ts = ts.QuadJoin(s)
	lts := len(ts)
	if lts != 128 {
		t.Fail()
	}
}

// Function TestLinearTranslation takes a 128-bit Bitstring and performs
// a translation then an inverse translation which should return the
// original Bitstring
func TestLinearTranslation(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestLinearTranslation\n")
	LtBitstring := LT(bs)
	LtiBitstring := LTInverse(LtBitstring)
	if LtiBitstring != bs {
		t.Fail()
	}
}

// Function TestSHat applies the parallel array of 32 copies of S-Box #3 and
// then applies the inverse to return the Bitstring to it's original form.
func TestSHat(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestSHat\n")
	bs1 := SHat(3, bs)
	orig := SHatInverse(3, bs1)
	if orig != bs {
		t.Fail()
	}
}

// Function TestLTBitslice applies the equations-based linear transformation
// and then reverses it.
func TestLTBitslice(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestLTBitslice\n")
	bss1 := bs.QuadSplit()
	bsl1 := LTBitslice(bss1)
	bss2 := LTBitsliceInverse(bsl1)
	for i := 0; i < 4; i++ {
		if bss2[i] != bss1[i] {
			t.Fail()
		}
	}
}

// Function TestPTable applies the permutation table sequence.
func TestPTable(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestPTable\n")
	bspi := IP(bs)
	bspf := FP(bspi)
	bspfr := FPInverse(bspf)
	bspir := IPInverse(bspfr)
	if bspir != bs {
		t.Fail()
	}
}

// Function TestS checks S function
func TestS(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestS\n")
	var target Bitstring = "0111"
	result := S(2, "0101")
	if result != target {
		t.Errorf("Output from S does not match target\n")
		t.Fail()
	}
}

// Function TestKeygen checks the formatting of the userkey.
func TestKeygen(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestKeygen\n")
	// Expected output from makeLongkey from 128-bit Bitstring 'bs'
	var longkey Bitstring = "110011100101001100101001010100101011100010" +
		"0110010100101001010111001010010101011001001001010110101001" +
		"0100101010010101001010100101100000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000"
	// Expected output in K[3]
	var k3 Bitstring = "010001101101011101010100110100011111110001011001" +
		"00000011100010100011001010110000000101111110001000000011001" +
		"011011001011010111011"
	// Expected KHat[3] from makeSubkeys using above longkey
	var khat3 Bitstring = "010011000110011001001100101100011010110000111" +
		"1100101100110001101000110000000101100001011011101101111101" +
		"0001110010101000001111001"
	var K Bitslice
	var KHat Bitslice
	ukey := bs
	ukeyl := makeLongkey(ukey)
	ukeyll := len(ukeyl)
	if ukeyll == 256 {
		if ukeyl != longkey {
			t.Errorf("Long key does not match expected\n")
			t.Fail()
		}
		K, KHat = makeSubkeys(ukeyl)
	}
	if len(K) != 33 && len(KHat) != 33 {
		t.Fail()
	}
	if K[3] != k3 {
		t.Errorf("K[3] does not match expected output\n")
		t.Fail()
	}
	if KHat[3] != khat3 {
		t.Errorf("KHat[3] does not match expected output\n")
		t.Fail()
	}
}

// Function TestR checks the R function returns the correct value
func TestR(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestR\n")
	var target Bitstring
	target = "00111011110100111000111011100000111101011001100000000001" +
		"0010100100111100010100110011010010001110011000000000011000" +
		"00110101110010"
	_, KHat := makeSubkeys(makeLongkey(bs))
	output := R(2, bs, KHat)
	if output != target {
		t.Errorf("Output doesn't match target\n")
		t.Fail()
	}
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
