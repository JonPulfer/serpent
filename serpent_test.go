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

// Function TestRotateLeft checks the RotateLeft method.
func TestRotateLeft(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestRotateLeft\n")
	var target Bitstring = "01010100101110011100101001100101001010100101" +
		"01110001001100101001010010101110010100101010110010010010101" +
		"1010100101001010100101010"
	bsLeft := bs.RotateLeft(11)
	if bsLeft != target {
		t.Errorf("bsLeft does not match target\n")
		t.Fail()
	}
}

// Function TestXor checks the Xor function.
func TestXor(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestXor\n")
	var target Bitstring = "0000"
	output := target.Xor(Bitslice{"0000", "0000"})
	if output != target {
		t.Errorf("Xor of two 0 Bitstrings does not equal" +
			" 0 Bitstring\n")
		t.Fail()
	}
	target = "1001"
	output = target.Xor(Bitslice{"1000", "0001"})
	if output != target {
		t.Errorf("Xor of two Bitstrings with 1 at either end" +
			" does not result in Bitstring with 1 at both ends\n")
		t.Fail()
	}
	target = "00100"
	output = target.Xor(Bitslice{"00001", "00100", "00001"})
	if output != target {
		t.Errorf("Xor of multiple Bitstrings does not match expected" +
			" output\n")
		t.Fail()
	}
}

// Function TestLinearTranslation takes a 128-bit Bitstring and performs
// a translation then an inverse translation which should return the
// original Bitstring
func TestLinearTranslation(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestLinearTranslation\n")
	var ltTarget Bitstring = "010100001000001110000000010110111011111110" +
		"11011100000110011010001110010100110100001011100101011011000" +
		"110011110000000011011101110"
	LtBitstring := LT(bs)
	if LtBitstring != ltTarget {
		t.Errorf("LT does not match target\n")
		t.Fail()
	}
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
	var bss1Target Bitslice = Bitslice{"00000001010101011010100101000110",
		"11011011100100101111100011110010",
		"10110001001110011100011111111100",
		"10111011101100011010001010101011"}
	bss1 := bs.QuadSplit()
	bsl1 := LTBitslice(bss1)
	for i := 0; i < 4; i++ {
		if bss1[i] != bss1Target[i] {
			t.Errorf("bss1[%d] is not correct\n", i)
			t.Fail()
		}
	}
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
	var target Bitstring = "11000100111101000100010110110010111000100000" +
		"10110011101001001001010110011000010110010001001001100011011" +
		"0101100110110011010011001"
	testFP := FP(bs)
	if testFP != target {
		t.Errorf("FP does not match target\n")
		t.Fail()
	}
	testFPInverse := FPInverse(target)
	if testFPInverse != bs {
		t.Errorf("FPInverse does not yield bs\n")
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

// Function TestSBitslice checks the SBitslice functions
func TestSBitslice(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestSBitslice\n")
	var bsnew Bitstring
	bslice := SBitslice(2, bs.QuadSplit())
	bssplit := SBitsliceInverse(2, bslice)
	bsnew = bsnew.QuadJoin(bssplit)
	if bsnew != bs {
		t.Errorf("bs not decoded correctly\n")
		t.Fail()
	}
}

// Function TestSBoxBitstring checks whether the loaded SBoxBitstring structure
// provides the correct pattern / value relationships.
func TestSBoxBitstring(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestSBoxBitstring\n")
	normal := SBoxBitstring[4]["1110"]
	if normal != "0110" {
		t.Errorf("Normal SBoxBitstring result is not expected value\n")
		t.Fail()
	}
	inverse := SBoxBitstringInverse[4]["0110"]
	if inverse != "1110" {
		t.Errorf("SBoxBitstringInverse result is not expected value\n")
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
	// Expected K[13]
	var k13 Bitstring = "11010100111111001101000010100101110011000000011" +
		"00010110100110010101000011111111010101100110011010111101101" +
		"0111110010010011101011"
	// Expected KHat[3] from makeSubkeys using above longkey
	var khat3 Bitstring = "010011000110011001001100101100011010110000111" +
		"1100101100110001101000110000000101100001011011101101111101" +
		"0001110010101000001111001"
	// Expected KHat[15]
	var khat15 Bitstring = "10110100110111001000100111011110010110001010" +
		"01011101101011110001110101101010000001001101010110100110001" +
		"1000100001011101000110101"
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
	if K[13] != k13 {
		t.Errorf("K[13] does not match expected output\n")
		t.Fail()
	}
	if KHat[3] != khat3 {
		t.Errorf("KHat[3] does not match expected output\n")
		t.Fail()
	}
	if KHat[15] != khat15 {
		t.Errorf("KHat[15] does not match expected output\n")
		t.Fail()
	}
}

// Function TestR checks the R function returns the correct value
func TestR(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestR\n")
	var target Bitstring = "00010011010001101100000010110000100011100110" +
		"11000010011010011101011101000111001010101100111001000101000" +
		"1110011101011110010101110"
	_, KHat := makeSubkeys(makeLongkey(bs))
	output := bs
	for i := 0; i < round; i++ {
		output = R(i, output, KHat)
	}
	if output != target {
		t.Errorf("Output doesn't match target\n")
		t.Fail()
	}
}

// Function TestRInverse checks the RInverse funtions returns bs
func TestRInverse(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestRInverse\n")
	var start Bitstring = "00010011010001101100000010110000100011100110" +
		"11000010011010011101011101000111001010101100111001000101000" +
		"1110011101011110010101110"
	_, KHat := makeSubkeys(makeLongkey(bs))
	var BHati Bitstring = start
	for i := round - 1; i >= 0; i-- {
		BHati = RInverse(i, BHati, KHat)
	}
	if BHati != bs {
		t.Errorf("Output from RInverse does not match bs\n")
		t.Fail()
	}
}

// Function TestRBitslice checks the RBitslice and RBitsliceInverse functions
// work correctly.
func TestRBitslice(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestRBitslice\n")
	var BiPlus1Target Bitstring = "1111100000110011100100011000010001111" +
		"00001010010111001100110010100001101001001000001100011111101" +
		"10111101101110010110111010111010"
	K, _ := makeSubkeys(makeLongkey(bs))
	BiPlus1 := RBitslice(2, bs, K)
	if BiPlus1 != BiPlus1Target {
		t.Errorf("BiPlus1 is not correct\n")
		t.Fail()
	}
	Bi := RBitsliceInverse(2, BiPlus1, K)
	if Bi != bs {
		t.Errorf("BiPlus1 has not been decoded correctly\n")
		t.Fail()
	}
}

// Function TestEncrypt checks the normal encryption algorithm
func TestEncrypt(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestEncrypt\n")
	var plainText Bitstring = "10011010011010001101110001110100101001010" +
		"10010101001110001010010100101000000011111110100101111110000" +
		"0110101001110011001010110010"
	var target Bitstring = "11111101110001101000110011011111010011000011" +
		"11010101101101101100110101000100010011001110100101101101011" +
		"1001011010111011110100101"
	encrypted := Encrypt(plainText, makeLongkey(bs))
	if encrypted != target {
		t.Errorf("Encrypted does not match target\n")
		t.Fail()
	}
}

// Function TestDecrypt checks the normal decryption algorithm
func TestDecrypt(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestDecrypt\n")
	var plainText Bitstring = "10011010011010001101110001110100101001010" +
		"10010101001110001010010100101000000011111110100101111110000" +
		"0110101001110011001010110010"
	var cipher Bitstring = "11111101110001101000110011011111010011000011" +
		"11010101101101101100110101000100010011001110100101101101011" +
		"1001011010111011110100101"
	decrypted := Decrypt(cipher, makeLongkey(bs))
	if decrypted != plainText {
		t.Errorf("Decrypted does not match plainText\n")
		t.Fail()
	}
}

// Function TestEncryptBitslice checks the bitslice encryption algorithm
func TestEncryptBitslice(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestEncryptBitslice\n")
	var plainText Bitstring = "10011010011010001101110001110100101001010" +
		"10010101001110001010010100101000000011111110100101111110000" +
		"0110101001110011001010110010"
	var cipherTextBitslice Bitstring = "11111101110001101000110011011111" +
		"01001100001111010101101101101100110101000100010011001110100" +
		"1011011010111001011010111011110100101"
	cText := EncryptBitslice(plainText, makeLongkey(bs))
	if cText != cipherTextBitslice {
		t.Errorf("EncryptBitslice does not match cipherTextBitslice\n")
		t.Fail()
	}
}

// Function TestDecryptBitslice checks the bitslice decryption algorithm
func TestDecryptBitslice(t *testing.T) {
	fmt.Printf("Running\t-\t\tTestDecryptBitslice\n")
	var plainText Bitstring = "10011010011010001101110001110100101001010" +
		"10010101001110001010010100101000000011111110100101111110000" +
		"0110101001110011001010110010"
	var cipherTextBitslice Bitstring = "11111101110001101000110011011111" +
		"01001100001111010101101101101100110101000100010011001110100" +
		"1011011010111001011010111011110100101"
	pText := DecryptBitslice(cipherTextBitslice, makeLongkey(bs))
	if pText != plainText {
		t.Errorf("DecryptBitslice does not match plainText\n")
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
