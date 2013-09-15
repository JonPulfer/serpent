/*
Package serpent implements the serpent encryption cipher.

This implementation has been created based on Frank Stajano's (of
Cambridge University Computer Laboratory <http://www.cl.cam.ac.uk/~fms27/>)
python implementation

Serpent cipher invented by Ross Anderson, Eli Biham, Lars Knudsen.
*/
package main

import (
	"fmt"
)

type SBox []byte
type Ttable []int
type bitstring string

// Functions used in the formal description of the cipher

// Function S applies S-Box number 'box' to 4-bit bitstring 'input'
// and return a 4-bit bitstring.
func S(box SBox, input byte) byte {
	return SBoxBitstring[box%8][input]
}

// Data tables

// Each element of this list corresponds to one S-Box. Each S-Box in turn is
// a list of 16 integers in the range 0..15, without repetitions. Having the
// value v (say, 14) in position p (say, 0) means that if the input to that
// S-Box is the pattern p (0, or 0x0) then the output will the pattern v
// (14, or 0xe).
var SBoxDecimalTable []SBox = []SBox{
	[]byte{3, 8, 15, 1, 10, 6, 5, 11, 14, 13, 4, 2, 7, 0, 9, 12}, // S0
	[]byte{15, 12, 2, 7, 9, 0, 5, 10, 1, 11, 14, 8, 6, 13, 3, 4}, // S1
	[]byte{8, 6, 7, 9, 3, 12, 10, 15, 13, 1, 14, 4, 0, 11, 5, 2}, // S2
	[]byte{0, 15, 11, 8, 12, 9, 6, 3, 13, 1, 2, 4, 10, 7, 5, 14}, // S3
	[]byte{1, 15, 8, 3, 12, 0, 11, 6, 2, 5, 4, 10, 9, 14, 7, 13}, // S4
	[]byte{15, 5, 2, 11, 4, 10, 9, 12, 0, 3, 14, 8, 13, 6, 7, 1}, // S5
	[]byte{7, 2, 12, 5, 8, 4, 6, 11, 14, 9, 1, 15, 13, 3, 10, 0}, // S6
	[]byte{1, 13, 15, 0, 14, 8, 2, 11, 7, 4, 12, 10, 9, 3, 5, 6}, // S7
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
var LTTable Ttable = []int{
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

func main() {
}
