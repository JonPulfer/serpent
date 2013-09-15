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

func main() {
}
