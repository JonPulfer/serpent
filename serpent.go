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
// S-Box is the pattern p (0, or 0x0) then the output will the pattern v
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

// The Initial and Final permutations are each represented by one list
// containing the integers in 0..127 without repetitions.  Having value v
// (say, 32) at position p (say, 1) means that the output bit at position p
// (1) comes from the input bit at position v (32).

var IPTable []int = []int{
	0, 32, 64, 96, 1, 33, 65, 97, 2, 34, 66, 98, 3, 35, 67, 99,
	4, 36, 68, 100, 5, 37, 69, 101, 6, 38, 70, 102, 7, 39, 71, 103,
	8, 40, 72, 104, 9, 41, 73, 105, 10, 42, 74, 106, 11, 43, 75, 107,
	12, 44, 76, 108, 13, 45, 77, 109, 14, 46, 78, 110, 15, 47, 79, 111,
	16, 48, 80, 112, 17, 49, 81, 113, 18, 50, 82, 114, 19, 51, 83, 115,
	20, 52, 84, 116, 21, 53, 85, 117, 22, 54, 86, 118, 23, 55, 87, 119,
	24, 56, 88, 120, 25, 57, 89, 121, 26, 58, 90, 122, 27, 59, 91, 123,
	28, 60, 92, 124, 29, 61, 93, 125, 30, 62, 94, 126, 31, 63, 95, 127,
}

var FPTable []int = []int{
	0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60,
	64, 68, 72, 76, 80, 84, 88, 92, 96, 100, 104, 108, 112, 116, 120, 124,
	1, 5, 9, 13, 17, 21, 25, 29, 33, 37, 41, 45, 49, 53, 57, 61,
	65, 69, 73, 77, 81, 85, 89, 93, 97, 101, 105, 109, 113, 117, 121, 125,
	2, 6, 10, 14, 18, 22, 26, 30, 34, 38, 42, 46, 50, 54, 58, 62,
	66, 70, 74, 78, 82, 86, 90, 94, 98, 102, 106, 110, 114, 118, 122, 126,
	3, 7, 11, 15, 19, 23, 27, 31, 35, 39, 43, 47, 51, 55, 59, 63,
	67, 71, 75, 79, 83, 87, 91, 95, 99, 103, 107, 111, 115, 119, 123, 127,
}

// The Linear Transformation is represented as a list of 128 lists, one for
// each output bit. Each one of the 128 lists is composed of a variable
// number of integers in 0..127 specifying the positions of the input bits
// that must be XORed together (say, 72, 144 and 125) to yield the output
// bit corresponding to the position of that list (say, 1).
var LTTable []Ttable = []Ttable{
	[]int{16, 52, 56, 70, 83, 94, 105},
	[]int{72, 114, 125},
	[]int{2, 9, 15, 30, 76, 84, 126},
	[]int{36, 90, 103},
	[]int{20, 56, 60, 74, 87, 98, 109},
	[]int{1, 76, 118},
	[]int{2, 6, 13, 19, 34, 80, 88},
	[]int{40, 94, 107},
	[]int{24, 60, 64, 78, 91, 102, 113},
	[]int{5, 80, 122},
	[]int{6, 10, 17, 23, 38, 84, 92},
	[]int{44, 98, 111},
	[]int{28, 64, 68, 82, 95, 106, 117},
	[]int{9, 84, 126},
	[]int{10, 14, 21, 27, 42, 88, 96},
	[]int{48, 102, 115},
	[]int{32, 68, 72, 86, 99, 110, 121},
	[]int{2, 13, 88},
	[]int{14, 18, 25, 31, 46, 92, 100},
	[]int{52, 106, 119},
	[]int{36, 72, 76, 90, 103, 114, 125},
	[]int{6, 17, 92},
	[]int{18, 22, 29, 35, 50, 96, 104},
	[]int{56, 110, 123},
	[]int{1, 40, 76, 80, 94, 107, 118},
	[]int{10, 21, 96},
	[]int{22, 26, 33, 39, 54, 100, 108},
	[]int{60, 114, 127},
	[]int{5, 44, 80, 84, 98, 111, 122},
	[]int{14, 25, 100},
	[]int{26, 30, 37, 43, 58, 104, 112},
	[]int{3, 118},
	[]int{9, 48, 84, 88, 102, 115, 126},
	[]int{18, 29, 104},
	[]int{30, 34, 41, 47, 62, 108, 116},
	[]int{7, 122},
	[]int{2, 13, 52, 88, 92, 106, 119},
	[]int{22, 33, 108},
	[]int{34, 38, 45, 51, 66, 112, 120},
	[]int{11, 126},
	[]int{6, 17, 56, 92, 96, 110, 123},
	[]int{26, 37, 112},
	[]int{38, 42, 49, 55, 70, 116, 124},
	[]int{2, 15, 76},
	[]int{10, 21, 60, 96, 100, 114, 127},
	[]int{30, 41, 116},
	[]int{0, 42, 46, 53, 59, 74, 120},
	[]int{6, 19, 80},
	[]int{3, 14, 25, 100, 104, 118},
	[]int{34, 45, 120},
	[]int{4, 46, 50, 57, 63, 78, 124},
	[]int{10, 23, 84},
	[]int{7, 18, 29, 104, 108, 122},
	[]int{38, 49, 124},
	[]int{0, 8, 50, 54, 61, 67, 82},
	[]int{14, 27, 88},
	[]int{11, 22, 33, 108, 112, 126},
	[]int{0, 42, 53},
	[]int{4, 12, 54, 58, 65, 71, 86},
	[]int{18, 31, 92},
	[]int{2, 15, 26, 37, 76, 112, 116},
	[]int{4, 46, 57},
	[]int{8, 16, 58, 62, 69, 75, 90},
	[]int{22, 35, 96},
	[]int{6, 19, 30, 41, 80, 116, 120},
	[]int{8, 50, 61},
	[]int{12, 20, 62, 66, 73, 79, 94},
	[]int{26, 39, 100},
	[]int{10, 23, 34, 45, 84, 120, 124},
	[]int{12, 54, 65},
	[]int{16, 24, 66, 70, 77, 83, 98},
	[]int{30, 43, 104},
	[]int{0, 14, 27, 38, 49, 88, 124},
	[]int{16, 58, 69},
	[]int{20, 28, 70, 74, 81, 87, 102},
	[]int{34, 47, 108},
	[]int{0, 4, 18, 31, 42, 53, 92},
	[]int{20, 62, 73},
	[]int{24, 32, 74, 78, 85, 91, 106},
	[]int{38, 51, 112},
	[]int{4, 8, 22, 35, 46, 57, 96},
	[]int{24, 66, 77},
	[]int{28, 36, 78, 82, 89, 95, 110},
	[]int{42, 55, 116},
	[]int{8, 12, 26, 39, 50, 61, 100},
	[]int{28, 70, 81},
	[]int{32, 40, 82, 86, 93, 99, 114},
	[]int{46, 59, 120},
	[]int{12, 16, 30, 43, 54, 65, 104},
	[]int{32, 74, 85},
	[]int{36, 90, 103, 118},
	[]int{50, 63, 124},
	[]int{16, 20, 34, 47, 58, 69, 108},
	[]int{36, 78, 89},
	[]int{40, 94, 107, 122},
	[]int{0, 54, 67},
	[]int{20, 24, 38, 51, 62, 73, 112},
	[]int{40, 82, 93},
	[]int{44, 98, 111, 126},
	[]int{4, 58, 71},
	[]int{24, 28, 42, 55, 66, 77, 116},
	[]int{44, 86, 97},
	[]int{2, 48, 102, 115},
	[]int{8, 62, 75},
	[]int{28, 32, 46, 59, 70, 81, 120},
	[]int{48, 90, 101},
	[]int{6, 52, 106, 119},
	[]int{12, 66, 79},
	[]int{32, 36, 50, 63, 74, 85, 124},
	[]int{52, 94, 105},
	[]int{10, 56, 110, 123},
	[]int{16, 70, 83},
	[]int{0, 36, 40, 54, 67, 78, 89},
	[]int{56, 98, 109},
	[]int{14, 60, 114, 127},
	[]int{20, 74, 87},
	[]int{4, 40, 44, 58, 71, 82, 93},
	[]int{60, 102, 113},
	[]int{3, 18, 72, 114, 118, 125},
	[]int{24, 78, 91},
	[]int{8, 44, 48, 62, 75, 86, 97},
	[]int{64, 106, 117},
	[]int{1, 7, 22, 76, 118, 122},
	[]int{28, 82, 95},
	[]int{12, 48, 52, 66, 79, 90, 101},
	[]int{68, 110, 121},
	[]int{5, 11, 26, 80, 122, 126},
	[]int{32, 86, 99},
}

var LTTableInverse []Ttable = []Ttable{
	[]int{53, 55, 72},
	[]int{1, 5, 20, 90},
	[]int{15, 102},
	[]int{3, 31, 90},
	[]int{57, 59, 76},
	[]int{5, 9, 24, 94},
	[]int{19, 106},
	[]int{7, 35, 94},
	[]int{61, 63, 80},
	[]int{9, 13, 28, 98},
	[]int{23, 110},
	[]int{11, 39, 98},
	[]int{65, 67, 84},
	[]int{13, 17, 32, 102},
	[]int{27, 114},
	[]int{1, 3, 15, 20, 43, 102},
	[]int{69, 71, 88},
	[]int{17, 21, 36, 106},
	[]int{1, 31, 118},
	[]int{5, 7, 19, 24, 47, 106},
	[]int{73, 75, 92},
	[]int{21, 25, 40, 110},
	[]int{5, 35, 122},
	[]int{9, 11, 23, 28, 51, 110},
	[]int{77, 79, 96},
	[]int{25, 29, 44, 114},
	[]int{9, 39, 126},
	[]int{13, 15, 27, 32, 55, 114},
	[]int{81, 83, 100},
	[]int{1, 29, 33, 48, 118},
	[]int{2, 13, 43},
	[]int{1, 17, 19, 31, 36, 59, 118},
	[]int{85, 87, 104},
	[]int{5, 33, 37, 52, 122},
	[]int{6, 17, 47},
	[]int{5, 21, 23, 35, 40, 63, 122},
	[]int{89, 91, 108},
	[]int{9, 37, 41, 56, 126},
	[]int{10, 21, 51},
	[]int{9, 25, 27, 39, 44, 67, 126},
	[]int{93, 95, 112},
	[]int{2, 13, 41, 45, 60},
	[]int{14, 25, 55},
	[]int{2, 13, 29, 31, 43, 48, 71},
	[]int{97, 99, 116},
	[]int{6, 17, 45, 49, 64},
	[]int{18, 29, 59},
	[]int{6, 17, 33, 35, 47, 52, 75},
	[]int{101, 103, 120},
	[]int{10, 21, 49, 53, 68},
	[]int{22, 33, 63},
	[]int{10, 21, 37, 39, 51, 56, 79},
	[]int{105, 107, 124},
	[]int{14, 25, 53, 57, 72},
	[]int{26, 37, 67},
	[]int{14, 25, 41, 43, 55, 60, 83},
	[]int{0, 109, 111},
	[]int{18, 29, 57, 61, 76},
	[]int{30, 41, 71},
	[]int{18, 29, 45, 47, 59, 64, 87},
	[]int{4, 113, 115},
	[]int{22, 33, 61, 65, 80},
	[]int{34, 45, 75},
	[]int{22, 33, 49, 51, 63, 68, 91},
	[]int{8, 117, 119},
	[]int{26, 37, 65, 69, 84},
	[]int{38, 49, 79},
	[]int{26, 37, 53, 55, 67, 72, 95},
	[]int{12, 121, 123},
	[]int{30, 41, 69, 73, 88},
	[]int{42, 53, 83},
	[]int{30, 41, 57, 59, 71, 76, 99},
	[]int{16, 125, 127},
	[]int{34, 45, 73, 77, 92},
	[]int{46, 57, 87},
	[]int{34, 45, 61, 63, 75, 80, 103},
	[]int{1, 3, 20},
	[]int{38, 49, 77, 81, 96},
	[]int{50, 61, 91},
	[]int{38, 49, 65, 67, 79, 84, 107},
	[]int{5, 7, 24},
	[]int{42, 53, 81, 85, 100},
	[]int{54, 65, 95},
	[]int{42, 53, 69, 71, 83, 88, 111},
	[]int{9, 11, 28},
	[]int{46, 57, 85, 89, 104},
	[]int{58, 69, 99},
	[]int{46, 57, 73, 75, 87, 92, 115},
	[]int{13, 15, 32},
	[]int{50, 61, 89, 93, 108},
	[]int{62, 73, 103},
	[]int{50, 61, 77, 79, 91, 96, 119},
	[]int{17, 19, 36},
	[]int{54, 65, 93, 97, 112},
	[]int{66, 77, 107},
	[]int{54, 65, 81, 83, 95, 100, 123},
	[]int{21, 23, 40},
	[]int{58, 69, 97, 101, 116},
	[]int{70, 81, 111},
	[]int{58, 69, 85, 87, 99, 104, 127},
	[]int{25, 27, 44},
	[]int{62, 73, 101, 105, 120},
	[]int{74, 85, 115},
	[]int{3, 62, 73, 89, 91, 103, 108},
	[]int{29, 31, 48},
	[]int{66, 77, 105, 109, 124},
	[]int{78, 89, 119},
	[]int{7, 66, 77, 93, 95, 107, 112},
	[]int{33, 35, 52},
	[]int{0, 70, 81, 109, 113},
	[]int{82, 93, 123},
	[]int{11, 70, 81, 97, 99, 111, 116},
	[]int{37, 39, 56},
	[]int{4, 74, 85, 113, 117},
	[]int{86, 97, 127},
	[]int{15, 74, 85, 101, 103, 115, 120},
	[]int{41, 43, 60},
	[]int{8, 78, 89, 117, 121},
	[]int{3, 90},
	[]int{19, 78, 89, 105, 107, 119, 124},
	[]int{45, 47, 64},
	[]int{12, 82, 93, 121, 125},
	[]int{7, 94},
	[]int{0, 23, 82, 93, 109, 111, 123},
	[]int{49, 51, 68},
	[]int{1, 16, 86, 97, 125},
	[]int{11, 98},
	[]int{4, 27, 86, 97, 113, 115, 127},
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
func (s Bitstring) BinaryXor(s2 Bitstring) (result Bitstring) {
	if len(s) != len(s2) {
		fmt.Printf("cannot binaryXor bitstrings " +
			"of different lengths\n")
	}
	for i, b := range s {
		if string(b) == string(s2[i]) {
			result = result + "0"
		} else {
			result = result + "1"
		}
	}
	return
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

// Function LT applies the table based version of the linear transformation
// to the 128-bit Bitstring 'input' and returns a 128-bit Bitstring.
func LT(input Bitstring) Bitstring {
	if len(input) != 128 {
		fmt.Printf("input is not 128 bits long\n")
	}
	var result Bitstring
	t := len(LTTable)
	for i := 0; i < t; i++ {
		var outputBit Bitstring = "0"
		for _, j := range LTTable[i] {
			bsj := Bitstring(input[j])
			outputBit = outputBit.BinaryXor(bsj)
		}
		result = result + Bitstring(outputBit)
	}
	return result
}

// Function LTInverse applies the inverse of the table based version of
// the linear transformation to the 128-bit Bitstring 'output' and returns a
// 128-bit Bitstring.
func LTInverse(output Bitstring) Bitstring {
	if len(output) != 128 {
		fmt.Printf("output is not 128 bits long\n")
	}
	var result Bitstring
	t := len(LTTableInverse)
	for i := 0; i < t; i++ {
		var inputBit Bitstring = "0"
		for _, j := range LTTableInverse[i] {
			bsj := Bitstring(output[j])
			inputBit = inputBit.BinaryXor(bsj)
		}
		result = result + Bitstring(inputBit)
	}
	return result
}

// Function LTBitslice applies the equations-based version of the linear
// transformation to 'x', a list of 4 32-bit Bitstrings, least significant
// Bitstring first. Returns a list of 4 32-bit Bitstrings.
func LTBitslice(x Bitslice) Bitslice {
	x[0] = x[0].RotateLeft(13)
	x[2] = x[2].RotateLeft(3)
	x[1] = x[1].Xor(Bitslice{x[0], x[2]})
	x[3] = x[3].Xor(Bitslice{x[2], x[0].ShiftLeft(3)})
	x[1] = x[1].RotateLeft(1)
	x[3] = x[3].RotateLeft(7)
	x[0] = x[0].Xor(Bitslice{x[1], x[3]})
	x[2] = x[2].Xor(Bitslice{x[3], x[1].ShiftLeft(7)})
	x[0] = x[0].RotateLeft(5)
	x[2] = x[2].RotateLeft(22)

	return x
}

// Function LTBitsliceInverse applies, in reverse, the equations-based
// version of the linear transformation to 'x', a list of 4 32-bit Bitstrings,
// least significant bit first. Returns a list of 4 32-bit Bitstrings.
func LTBitsliceInverse(x Bitslice) Bitslice {
	x[2] = x[2].RotateRight(22)
	x[0] = x[0].RotateRight(5)
	x[2] = x[2].Xor(Bitslice{x[3], x[1].ShiftLeft(7)})
	x[0] = x[0].Xor(Bitslice{x[1], x[3]})
	x[3] = x[3].RotateRight(7)
	x[1] = x[1].RotateRight(1)
	x[3] = x[3].Xor(Bitslice{x[2], x[0].ShiftLeft(3)})
	x[1] = x[1].Xor(Bitslice{x[0], x[2]})
	x[2] = x[2].RotateRight(3)
	x[0] = x[0].RotateRight(13)

	return x
}

// Function applyPermutation applies the permutation specified by the
// 128-element list 'ptable' to the 128-bit Bitstring 'input' and return
// a 128-bit Bitstring.
func applyPermutation(ptable []int, input Bitstring) Bitstring {
	ptlen := len(ptable)
	iplen := len(input)
	if iplen != ptlen {
		fmt.Printf("Input size (%d) doesn't match ptable size "+
			"(%d)\n", iplen, ptlen)
	}
	var result Bitstring
	for i := 0; i < ptlen; i++ {
		r := Bitstring(input[ptable[i]])
		result = result + r
	}

	return result
}

// Function IP applies the initial permutation table to the 128-bit
// Bitstring 'input' and returns the result.
func IP(input Bitstring) Bitstring {
	return applyPermutation(IPTable, input)
}

// Function FP applies the final permutation table to the 128-bit Bitstring
// 'input' and returns the result.
func FP(input Bitstring) Bitstring {
	return applyPermutation(FPTable, input)
}

// Function FPInverse applies the final permutation in reverse.
func FPInverse(output Bitstring) Bitstring {
	return FP(output)
}

// Function IPInverse applies the initial permutation in reverse.
func IPInverse(output Bitstring) Bitstring {
	return IP(output)
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

	if 0 <= i && i <= round - 2 {
		BHatiPlus1 = LT(SHati)
	} else if i == round - 1 {
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
	
	if 0 <= i && i <= round - 2 {
		SHati = LTInverse(BHatiPlus1)
	} else if i == round - 1 {
		SHati = xored.Xor(Bitslice{BHatiPlus1, KHat[round]})
	} else {
		fmt.Printf("Round is out of range\n")
	}

	xored = SHatInverse(i, SHati)
	BHati = xored.Xor(Bitslice{xored, KHat[i]})

	return BHati
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
