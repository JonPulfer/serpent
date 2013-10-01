package serpent

import (
	"fmt"
	"testing"
)

var bs Bitstring
func init() {
	bs = "11001110010100110010100101010010101110001001100101" +
	"00101001010111001010010101011001001001010110101001010010101001" +
	"0101001010100101"
}

// Function TestQuadSplit tests splitting a 128-bit Bitstring
// into a 4 element []Bitstring each of 32-bits.
func TestQuadSplit(t *testing.T) {
	s := bs.QuadSplit()
	sl := len(s)
	if sl == 4 {
		for i := 0; i < sl; i++ {
			sls := len(s[i])
			if sls != 32 {
				t.FailNow()
			}
		}
	}
}

// Function TestQuadJoin tests joining a 4 element []Bitstring of 32-bits
// each into a single 128-bit Bitstring.
func TestQuadJoin(t *testing.T) {
	var ts Bitstring
	s := bs.QuadSplit()
	ts = ts.QuadJoin(s)
	lts := len(ts)
	if lts != 128 {
		t.FailNow()
	}
}

// This example takes a bitstring of "000111" and shifts it left by 2 places 
// showing the visual effect of the little-endian representation.
func ExampleBitstring_ShiftLeft() {
	var bs1 Bitstring
	bs1 = "000111"
	output := bs1.ShiftLeft(2)
	fmt.Printf("%v\n", output)
	// Output:
	// 000001
}

// This example takes a bitstring of "000111" and shifts it right by 2 places
// showing the visual effect of the little-endian representation.
func ExampleBitstring_ShiftRight() {
	var bs1 Bitstring
	bs1 = "000111"
	output := bs1.ShiftRight(2)
	fmt.Printf("%v\n", output)
	// Output:
	// 011100
}
