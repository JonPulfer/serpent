package serpent

import (
	"fmt"
)

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
	return IP(output)
}

// Function IPInverse applies the initial permutation in reverse.
func IPInverse(output Bitstring) Bitstring {
	return FP(output)
}
