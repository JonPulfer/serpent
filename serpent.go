/*
Package serpent implements the serpent encryption cipher.

This implementation has been created based on Frank Stajano's (of
Cambridge University Computer Laboratory <http://www.cl.cam.ac.uk/~fms27/>)
python implementation

Serpent cipher invented by Ross Anderson, Eli Biham, Lars Knudsen.
*/
package serpent

import (
	"fmt"
)

type SBox []int
type Ttable []int
type Bitstring string
type Bitslice []Bitstring
type Bitmap map[int]Bitstring
type Hexstring string

// Data tables

// Each element of this list corresponds to one S-Box. Each S-Box in turn is
// a list of 16 integers in the range 0..15, without repetitions. Having the
// value v (say, 14) in position p (say, 0) means that if the input to that
// S-Box is the pattern p (0, or 0x0) then the output will be the pattern v
// (14, or 0xe).
var SBoxDecimalTable []SBox = []SBox{
	[]int{3, 8, 15, 1, 10, 6, 5, 11, 14, 13, 4, 2, 7, 0, 9, 12}, // S0
	[]int{15, 12, 2, 7, 9, 0, 5, 10, 1, 11, 14, 8, 6, 13, 3, 4}, // S1
	[]int{8, 6, 7, 9, 3, 12, 10, 15, 13, 1, 14, 4, 0, 11, 5, 2}, // S2
	[]int{0, 15, 11, 8, 12, 9, 6, 3, 13, 1, 2, 4, 10, 7, 5, 14}, // S3
	[]int{1, 15, 8, 3, 12, 0, 11, 6, 2, 5, 4, 10, 9, 14, 7, 13}, // S4
	[]int{15, 5, 2, 11, 4, 10, 9, 12, 0, 3, 14, 8, 13, 6, 7, 1}, // S5
	[]int{7, 2, 12, 5, 8, 4, 6, 11, 14, 9, 1, 15, 13, 3, 10, 0}, // S6
	[]int{1, 13, 15, 0, 14, 8, 2, 11, 7, 4, 12, 10, 9, 3, 5, 6}, // S7
}

// Constants
const (
	phi   int = 0x9e3779b9
	round int = 32
)

var SBoxBitstring []map[Bitstring]Bitstring
var SBoxBitstringInverse []map[Bitstring]Bitstring

// Initialise variables when this package is imported.
func init() {
	var bs Bitstring
	for _, sbox := range SBoxDecimalTable {
		var dict map[Bitstring]Bitstring = make(
			map[Bitstring]Bitstring, len(sbox))
		var inverseDict map[Bitstring]Bitstring = make(
			map[Bitstring]Bitstring, len(sbox))

		for boxindex, box := range sbox {
			index := bs.FromInt(boxindex, 4)
			value := bs.FromInt(box, 4)
			dict[index] = value
			inverseDict[value] = index
		}
		SBoxBitstring = append(SBoxBitstring, dict)
		SBoxBitstringInverse = append(SBoxBitstringInverse, inverseDict)
	}
}

// Methods for Bitslice

// Method Reverse returns the Bitslice with the elements reversed.
func (bs Bitslice) Reverse() []Bitstring {
	bslen := len(bs)
	new := make([]Bitstring, bslen)
	for i := 0; i < bslen; i++ {
		new[i] = Bitstring(bs[bslen-1-i])
	}
	return new
}

// Methods for Bitstring

// Translate n from integer to bitstring, padding it with 0s as
// necessary to reach the minimum length 'minlen'. 'n' must be >= 0 since
// the bitstring format is undefined for negative integers.
//
// Note that, while the bitstring format can represent arbitrarily large
// numbers, this is not so for Go's normal integer type: on a 32-bit machine,
// values of n >= 2^31 need to be expressed as int64 or
// they will "look" negative and won't work.
func (s Bitstring) FromInt(n int, l int) (result Bitstring) {
	if l < 1 {
		fmt.Printf("a bitstring must have a least 1 char\n")
	}
	if n < 0 {
		fmt.Printf("bitstring representation undefined for " +
			"negative numbers\n")
	}
	for n > 0 {
		if n&1 == 1 {
			result = result + "1"
		} else {
			result = result + "0"
		}
		n = n >> 1
	}
	for len(result) < l {
		result = result + "0"
	}
	return
}

// ByteSlice returns a []byte representation of the bitstring
func (s Bitstring) ByteSlice() (result []byte) {
	for _, char := range s {
		if string(char) == "0" {
			result = append(result, '0')
		} else {
			result = append(result, '1')
		}
	}
	return
}

// ToHex returns a 1-char hexstring of a 4 char bitstring
func (s Bitstring) ToHex() (h Hexstring) {
	if len(s) > 4 {
		fmt.Printf("Bitstring is more than 4 chars, " +
			"cannot be converted to hex char\n")
	}
	var bin2hex = map[Bitstring]Hexstring{
		"0000": "0", "1000": "1", "0100": "2", "1100": "3",
		"0010": "4", "1010": "5", "0110": "6", "1110": "7",
		"0001": "8", "1001": "9", "0101": "a", "1101": "b",
		"0011": "c", "1011": "d", "0111": "e", "1111": "f",
	}
	return bin2hex[s]
}

// FromHex returns a 4-char bitstring of a 1-char hexstring
func (s Bitstring) FromHex(h Hexstring) Bitstring {
	if len(h) > 1 {
		fmt.Printf("Hex string is more than 1 char, " +
			"cannot be converted to bitstring\n")
	}
	var hex2bin = map[Hexstring]Bitstring{
		"0": "0000", "1": "1000", "2": "0100", "3": "1100",
		"4": "0010", "5": "1010", "6": "0110", "7": "1110",
		"8": "0001", "9": "1001", "a": "0101", "b": "1101",
		"c": "0011", "d": "1011", "e": "0111", "f": "1111",
	}
	return hex2bin[h]
}

// ToHexstring returns the hexstring representation of the
// bitstring
func (s Bitstring) ToHexstring() (result Hexstring) {
	ln := len(s)
	var b Bitslice
	var tb []byte = make([]byte, 4)

	for i := 0; i < ln; i = i + 4 {
		tb[0] = s[i]
		tb[1] = s[i+1]
		tb[2] = s[i+2]
		tb[3] = s[i+3]
		b = append(b, Bitstring(tb))
	}
	for _, nbs := range b {
		result = nbs.ToHex() + result
	}
	return
}

// ToBistring returns the bitstring representation of the
// hexstring
func (h Hexstring) ToBitstring() (result Bitstring) {
	ln := len(h)
	var rh []byte = make([]byte, ln)
	var n int = 0
	for i := ln - 1; i >= 0; i-- {
		rh[i] = h[n]
		n++
	}
	for j := 0; j < ln; j++ {
		result = result + result.FromHex(Hexstring(rh[j]))
	}
	return
}

// Return the xor of two bitstrings of equal length as another
// bitstring of the same length.
func (s Bitstring) BinaryXor(s2 Bitstring) Bitstring {
	if len(s) != len(s2) {
		fmt.Printf("cannot binaryXor bitstrings " +
			"of different lengths\n")
	}
	var result Bitstring = ""
	for i, b := range s {
		if uint8(b) == s2[i] {
			result = result + "0"
		} else {
			result = result + "1"
		}
	}
	return result
}

// Return the xor of an arbitrary number of bitstrings of the same
// length as another bitstring of the same length.
func (s Bitstring) Xor(args Bitslice) (result Bitstring) {
	if len(args) == 0 {
		fmt.Printf("at least one argument needed\n")
	}
	result = args[0]
	for _, arg := range args[1:] {
		result = result.BinaryXor(arg)
	}
	return
}

// Take a bitstring 'input' of arbitrary length. Rotate it left by
// 'places' places. Left means that the 'places' most significant bits are
// taken out and reinserted as the least significant bits. Note that,
// because the bitstring representation is little-endian, the visual
// effect is actually that of rotating the string to the right.
func (input Bitstring) RotateLeft(places int) Bitstring {
	wc := input.ByteSlice()
	lw := len(wc)
	var nc []byte = make([]byte, lw)
	var op int
	for i := 0; i < lw; i++ {
		if i < places {
			op = lw - places + i
		} else if i == places {
			op = 0
		} else if i > places {
			op = i - places
		}
		nc[i] = wc[op]
	}
	return Bitstring(nc)
}

// Take a bitstring 'input' of arbitrary length and rotate it right
// by 'places' places.
func (input Bitstring) RotateRight(places int) Bitstring {
	wc := input.ByteSlice()
	lw := len(wc)
	var nc []byte = make([]byte, lw)
	var op int
	for i := 0; i < lw; i++ {
		if i+places < lw {
			op = i + places
		} else if i+places == lw {
			op = 0
		} else if i+places > lw {
			op = places + i - lw
		}
		nc[i] = wc[op]
	}
	return Bitstring(nc)
}

// Take a bitstring 's' of arbitrary length. Shift it left by 'p'
// places. Left means that the 'p' most significant bits are shifted out
// and dropped, while 'p' 0s are inserted in the the least significant
// bits. Note that, because the bitstring representation is little-endian,
// the visual effect is actually that of shifting the string to the
// right. Negative values for 'p' are allowed, with the effect of shifting
// right instead (i.e. the 0s are inserted in the most significant bits).
func (s Bitstring) ShiftLeft(places int) Bitstring {
	wc := s.ByteSlice()
	lw := len(wc)
	var nc []byte = make([]byte, lw)
	if places < 0 {
		return s.ShiftRight(places - places*2)
	}
	for i := 0; i < lw; i++ {
		if i < places {
			nc[i] = '0'
		} else if i == places {
			nc[places] = wc[i-places]
		} else {
			nc[i] = wc[i-places]
		}
	}
	return Bitstring(nc)
}

// Take a bitstring 's' of arbitrary length and shift it right. Same
// as Bitstring.ShiftLeft using negative int.
func (s Bitstring) ShiftRight(places int) Bitstring {
	wc := s.ByteSlice()
	lw := len(wc)
	var nc []byte = make([]byte, lw)
	for i := 0; i < lw; i++ {
		if i <= places {
			nc[i] = wc[places+i]
		} else if i+places < lw {
			nc[i] = wc[i+places]
		} else {
			nc[i] = '0'
		}
	}
	return Bitstring(nc)
}

// QuadSplit breaks a 128-bit bitstring into 4 32-bit bitstrings
// and returns them, least significant bitstring first
func (s Bitstring) QuadSplit() Bitslice {
	if len(s) != 128 {
		fmt.Printf("Bitstring must be 128-bits to be quadsplit\n")
	}
	result := make(Bitslice, 4)

	for i := 0; i < 4; i++ {
		result[i] = s[i*32 : (i+1)*32]
	}
	return result
}

// QuadJoin joins 4 32-bit bitstrings into a single 128-bit
// bitstring.
func (s Bitstring) QuadJoin(bs Bitslice) Bitstring {
	if len(bs) != 4 {
		fmt.Printf("List of bitstrings must " +
			"contain 4 * 32-bit bitstrings\n")
	}
	return bs[0] + bs[1] + bs[2] + bs[3]
}

// Functions used in the formal description of the cipher

// Function S applies S-Box number 'box' to 4-bit bitstring 'input'
// and return a 4-bit bitstring.
func S(box int, input Bitstring) Bitstring {
	return SBoxBitstring[box%8][input]
}

// Function SInverse applies S-Box number box in reverse to 4-bit bitstring
// 'output' and return a 4-bit bitstring 'input' as the result
func SInverse(box int, output Bitstring) Bitstring {
	return SBoxBitstringInverse[box%8][output]
}

// Function SHat applies a parallel array of 32 copies of S-Box number 'box'
// to the 128-bit bitstring 'input' and return a 128-bit bitstring as the
// result
func SHat(box int, input Bitstring) Bitstring {
	var bs Bitstring
	for i := 0; i < 32; i++ {
		bs = bs + S(box, input[4*i:4*(i+1)])
	}
	return bs
}

// Function SHatInverse applies in reverse, a parallel array of 32 copies of
// S-Box number 'box' to the 128-bit bitstring 'output' and return a 128-bit
// bitstring (the input) as the result
func SHatInverse(box int, output Bitstring) Bitstring {
	var bs Bitstring
	for i := 0; i < 32; i++ {
		bs = bs + SInverse(box, output[4*i:4*(i+1)])
	}
	return bs
}

// Function SBitslice takes 'words', a list of 4 32-bit bitstring, least
// significant word first and returns a similar list of 4 32-bit bitstrings.
// Obtained as follows: -
//
// For each bit position from 0 to 31, apply S-Box number 'box' to the 4 input
// bits coming from the current position in each of the items in 'words' and
// put the 4 output bits in the corresponding positions in the output
// words.
func SBitslice(box int, words Bitslice) Bitslice {
	result := make(Bitslice, 4)
	for i := 0; i < 32; i++ {
		var c0 Bitstring = Bitstring(int(words[0][i]))
		var c1 Bitstring = Bitstring(int(words[1][i]))
		var c2 Bitstring = Bitstring(int(words[2][i]))
		var c3 Bitstring = Bitstring(int(words[3][i]))
		quad := S(box, Bitstring(c0+c1+c2+c3))
		for j := 0; j < 4; j++ {
			result[j] = result[j] + Bitstring(int(quad[j]))
		}
	}
	return result
}

// Function SBitsliceInverse takes 'words', a list of 4 32-bit bitstring, least
// significant word first and returns a similar list of 4 32-bit bitstrings.
// Obtained as follows: -
//
// For each bit position from 0 to 31, apply S-Box number 'box' in reverse
// to the 4 input bits coming from the current position in each of the items
// in 'words' and put the 4 output bits in the corresponding positions in the
// output words.
func SBitsliceInverse(box int, words Bitslice) Bitslice {
	result := make(Bitslice, 4)
	for i := 0; i < 32; i++ {
		var c0 Bitstring = Bitstring(int(words[0][i]))
		var c1 Bitstring = Bitstring(int(words[1][i]))
		var c2 Bitstring = Bitstring(int(words[2][i]))
		var c3 Bitstring = Bitstring(int(words[3][i]))
		quad := SInverse(box, Bitstring(c0+c1+c2+c3))
		for j := 0; j < 4; j++ {
			result[j] = result[j] + Bitstring(int(quad[j]))
		}
	}
	return result
}

// Function R applies round 'i' to the 128-bit Bitstring 'BHati', returning
// another 128-bit Bitstring (conceptually BHatiPlus1). Do this using the
// appropriately numbered subkey(s) from the 'KHat' list of 33 128-bit
// Bitstrings.
func R(i int, BHati Bitstring, KHat Bitslice) Bitstring {
	var xored Bitstring
	var BHatiPlus1 Bitstring
	xored = xored.Xor(Bitslice{BHati, KHat[i]})
	SHati := SHat(i, xored)

	if 0 <= i && i <= round-2 {
		BHatiPlus1 = LT(SHati)
	} else if i == round-1 {
		BHatiPlus1 = BHatiPlus1.Xor(Bitslice{SHati, KHat[round]})
	} else {
		fmt.Printf("Round is out of range\n")
	}

	return BHatiPlus1
}

// Function RInverse applies the round 'i' in reverse to the 128-bit Bitstring
// 'BHatiPlus1', returning another 128-bit Bitstring (conceptually BHati). Do
// this using the appropriately numbered subkey(s) from the 'KHat' list of 33
// 128-bit Bitstrings.
func RInverse(i int, BHatiPlus1 Bitstring, KHat Bitslice) Bitstring {
	var xored Bitstring
	var BHati Bitstring
	var SHati Bitstring

	if 0 <= i && i <= round-2 {
		SHati = LTInverse(BHatiPlus1)
	} else if i == round-1 {
		SHati = xored.Xor(Bitslice{BHatiPlus1, KHat[round]})
	} else {
		fmt.Printf("Round is out of range\n")
	}

	xored = SHatInverse(i, SHati)
	BHati = xored.Xor(Bitslice{xored, KHat[i]})

	return BHati
}

// Function RBitslice applies round 'i' (Bitslice version) to the 128-bit
// Bitstring 'Bi' and return another 128-bit Bitstring (conceptually B i+1).
// Use the appropriately numbered subkey(s) from the 'K' list of 33 128-bit
// Bitstrings.
func RBitslice(i int, Bi Bitstring, K Bitslice) Bitstring {
	var xored Bitstring
	var BiPlus1 Bitstring

	// 1. Key mixing
	xored = xored.Xor(Bitslice{Bi, K[i]})

	// 2. S Boxes
	Si := SBitslice(i, xored.QuadSplit())

	// 3. Linear Transformation
	if i == round-1 {
		// In the last round, replaced by an additional key mixing
		BiPlus1 = xored.Xor(Bitslice{xored.QuadJoin(Si), K[round]})
	} else {
		BiPlus1 = xored.QuadJoin(LTBitslice(Si))
	}

	return BiPlus1
}

// Function RBitsliceInverse applies the inverse of round 'i' (bitslice
// version) to the 128-bit Bitstring 'BiPlus1' and return another 128-bit
// Bitstring (conceptually B i). Use the appropriately numbered subkey(s) from
// the 'K' list of 33 128-bit Bitstrings.
func RBitsliceInverse(i int, BiPlus1 Bitstring, K Bitslice) Bitstring {
	var xoredbitslice Bitslice
	var Bi Bitstring
	var SiTemp Bitstring
	var Si Bitslice

	// 3. Linear Transformation
	if i == round-1 {
		// In the last round, replaced by an additional key mixing
		SiTemp = SiTemp.Xor(Bitslice{BiPlus1, K[round]})
		Si = SiTemp.QuadSplit()
	} else {
		Si = LTBitsliceInverse(BiPlus1.QuadSplit())
	}

	// 2. S Boxes
	xoredbitslice = SBitsliceInverse(i, Si)

	// 1. Key mixing
	Bi = Bi.Xor(Bitslice{Bi.QuadJoin(xoredbitslice), K[i]})

	return Bi
}

// Function makeSubkeys takes the 256-bit Bitstring 'userkey' and returns two
// lists (conceptually K and KHat) of 33 128-bit Bitstrings each.
func makeSubkeys(userkey Bitstring) (Bitslice, Bitslice) {
	// Convert the userkey to 8 32-bit words.
	w := make(Bitmap, 132)
	for i := -8; i < 0; i++ {
		w[i] = userkey[(i+8)*32 : (i+9)*32]
	}

	// Expand the 8 words to a prekey w0 ... w131 with the affine
	// recurrence.
	var tempbs Bitstring
	for i := 0; i < 132; i++ {
		tempbsl := Bitslice{w[i-8], w[i-5], w[i-3], w[i-1],
			tempbs.FromInt(phi, 32),
			tempbs.FromInt(i, 32)}
		tempbs = tempbs.Xor(tempbsl)
		w[i] = tempbs.RotateLeft(11)
	}

	// The round keys are now calculated from the prekeys using the
	// S-Boxes in bitslice mode. Each k[i] is a 32-bit Bitstring.
	k := make(Bitslice, 132)
	for i := 0; i < round+1; i++ {
		whichS := (round + 3 - i) % round
		k[0+4*i] = ""
		k[1+4*i] = ""
		k[2+4*i] = ""
		k[3+4*i] = ""
		var input Bitstring
		for j := 0; j < 32; j++ {
			input = Bitstring(w[0+4*i][j]) +
				Bitstring(w[1+4*i][j]) +
				Bitstring(w[2+4*i][j]) +
				Bitstring(w[3+4*i][j])
			output := S(whichS, input)
			for l := 0; l < 4; l++ {
				k[l+4*i] = k[l+4*i] + Bitstring(output[l])
			}
		}
	}

	// We then renumber the 32-bit values k_j as 128-bit subkeys K_i
	K := Bitslice{}
	for i := 0; i < 33; i++ {
		K = append(K, k[4*i]+k[4*i+1]+k[4*i+2]+k[4*i+3])
	}

	// We now apply IP to the round key in order to place the key bits
	// in the correct column.
	KHat := Bitslice{}
	for i := 0; i < 33; i++ {
		KHat = append(KHat, IP(K[i]))
	}

	return K, KHat
}

// Function makeLongkey takes a bitstring key 'k' and returns the long
// (256-bit) version of that key.
func makeLongkey(k Bitstring) Bitstring {
	lk := len(k)
	if lk%32 != 0 || lk < 64 || lk > 256 {
		fmt.Printf("Invalid key length(%d bits)", lk)
	}
	if lk == 256 {
		return k
	}

	for i := 0; i < 256-lk; i++ {
		if i == 0 {
			k = k + "1"
		} else {
			k = k + "0"
		}
	}

	return k
}

// Function Encrypt uses the 256-bit Bitstring 'userKey' to encrypt the
// 128-bit Bitstring 'plainText' by the normal algorithm. Returns a 128-bit
// cipher text Bitstring.
func Encrypt(plainText Bitstring, userKey Bitstring) Bitstring {
	_, KHat := makeSubkeys(userKey)
	BHat := IP(plainText)
	for i := 0; i < round; i++ {
		BHat = R(i, BHat, KHat)
	}
	C := FP(BHat)

	return C
}

// Function EncryptBitslice encrypts the 128-bit Bitstring 'plainText' with
// the 256-bit Bitstring 'userKey' using the bitslice algorithm. Returns a
// 128-bit cipher text Bitstring.
func EncryptBitslice(plainText Bitstring, userKey Bitstring) Bitstring {
	K, _ := makeSubkeys(userKey)
	B := plainText
	for i := 0; i < round; i++ {
		B = RBitslice(i, B, K)
	}

	return B
}

// Function Decrypt uses the 256-bit Bitstring 'userKey' to decrypt the
// 128-bit Bitstring 'cipherText' using the normal algorithm. Returns a
// 128-bit Bitstring which is the plain text.
func Decrypt(cipherText Bitstring, userKey Bitstring) Bitstring {
	_, KHat := makeSubkeys(userKey)
	BHat := FPInverse(cipherText)
	for i := round - 1; i >= 0; i-- {
		BHat = RInverse(i, BHat, KHat)
	}
	plainText := IPInverse(BHat)

	return plainText
}

// Function DecryptBitslice decrypts the 128-bit Bitstring 'cipherText' with
// the 256-bit Bitstring 'userKey' using the bitslice algorithm. Returns a
// 128-bit Bitstring which is the plain text.
func DecryptBitslice(cipherText Bitstring, userKey Bitstring) Bitstring {
	K, _ := makeSubkeys(userKey)
	B := cipherText
	for i := round - 1; i >= 0; i-- {
		B = RBitsliceInverse(i, B, K)
	}

	return B
}
