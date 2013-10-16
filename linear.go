package serpent

import (
	"fmt"
)

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
	x[1] = x[1].Xor(Bitslice{x[1], x[0], x[2]})
	x[3] = x[3].Xor(Bitslice{x[3], x[2], x[0].ShiftLeft(3)})
	x[1] = x[1].RotateLeft(1)
	x[3] = x[3].RotateLeft(7)
	x[0] = x[0].Xor(Bitslice{x[0], x[1], x[3]})
	x[2] = x[2].Xor(Bitslice{x[2], x[3], x[1].ShiftLeft(7)})
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
	x[2] = x[2].Xor(Bitslice{x[2], x[3], x[1].ShiftLeft(7)})
	x[0] = x[0].Xor(Bitslice{x[0], x[1], x[3]})
	x[3] = x[3].RotateRight(7)
	x[1] = x[1].RotateRight(1)
	x[3] = x[3].Xor(Bitslice{x[3], x[2], x[0].ShiftLeft(3)})
	x[1] = x[1].Xor(Bitslice{x[1], x[0], x[2]})
	x[2] = x[2].RotateRight(3)
	x[0] = x[0].RotateRight(13)

	return x
}
