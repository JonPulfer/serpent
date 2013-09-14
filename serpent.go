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
	[]byte{ 3, 8,15, 1,10, 6, 5,11,14,13, 4, 2, 7, 0, 9,12 }, // S0
        []byte{15,12, 2, 7, 9, 0, 5,10, 1,11,14, 8, 6,13, 3, 4 }, // S1
        []byte{ 8, 6, 7, 9, 3,12,10,15,13, 1,14, 4, 0,11, 5, 2 }, // S2
        []byte{ 0,15,11, 8,12, 9, 6, 3,13, 1, 2, 4,10, 7, 5,14 }, // S3
        []byte{ 1,15, 8, 3,12, 0,11, 6, 2, 5, 4,10, 9,14, 7,13 }, // S4
        []byte{15, 5, 2,11, 4,10, 9,12, 0, 3,14, 8,13, 6, 7, 1 }, // S5
        []byte{ 7, 2,12, 5, 8, 4, 6,11,14, 9, 1,15,13, 3,10, 0 }, // S6
        []byte{ 1,13,15, 0,14, 8, 2,11, 7, 4,12,10, 9, 3, 5, 6 }, // S7
    }

func main() {
}
