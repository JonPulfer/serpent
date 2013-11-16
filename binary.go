package serpent

type uint32le [4]byte
type uint128le [16]byte
type uint256le [32]byte
type bitsliceNew []uint32le

const (
	digits = "0123456789abcdefghijklmnopqrstuvwxyz"
)

var shifts = [len(digits) + 1]uint{
	1 << 1: 1,
	1 << 2: 2,
	1 << 3: 3,
	1 << 4: 4,
	1 << 5: 5,
}

// Function newFromUint64 creates a uint128le from uint64 'num'.
func newFromUint64(num uint64) (result uint128le) {
	for i := 0; i < 8; i++ {
		result[i] = byte(num & 0xff)
		num >>= 8
	}
	return result
}

// Function newUint32le takes a []byte and returns a uint32le.
func newUint32le(num []byte) uint32le {
	var result uint32le
	for i, _ := range num {
		result[i] = num[i]
	}

	return result
}

// Function xor performs a sequential Xor operation on the elements
// in 'args'.
func xorNEW(args []uint128le) *uint128le {
	result := &args[0]
	for i := 1; i < len(args); i++ {
		result.binaryXorNEW(args[i])
	}
	return result
}

// Method quadSplit splits a 128-bit bitstring into a list of 4 32-bit
// bitstrings.
func (num uint128le) quadSplitNEW() []uint32le {
	result := bitsliceNew{newUint32le(num[0:3]),
		newUint32le(num[4:7]),
		newUint32le(num[8:11]),
		newUint32le(num[12:15])}
	return result
}

// Method binaryXor performs a Xor operation between 'num' and 'other'.
func (num *uint128le) binaryXorNEW(other uint128le) {
	for i := 0; i < 16; i++ {
		num[i] ^= other[i]
	}
}

// Function ROL rotates left by 'p' places.
func ROL(a uint128le, p uint) uint128le {
	for ; p > 0; p-- {
		var leftMostBit byte
		for i := 0; i < len(a); i++ {
			var leftBit byte
			if i == 0 {
				leftMostBit = a[i] & 0x80
			} else {
				leftBit = a[i] & 0x80
				a[i-1] += leftBit >> 7
			}
			a[i] <<= 1
		}
		a[15] += leftMostBit >> 7
	}
	return a
}

// Function ROR rotates right by 'p' places.
func ROR(a uint128le, p uint) uint128le {
	for ; p > 0; p-- {
		var rightMostBit byte
		for i := len(a) - 1; i >= 0; i-- {
			var rightBit byte
			if i == len(a)-1 {
				rightMostBit = a[i] & 0x1
			} else {
				rightBit = a[i] & 0x1
				a[i+1] += rightBit << 7
			}
			a[i] >>= 1
		}
		a[0] += rightMostBit << 7
	}
	return a
}

// Method shiftLeft shifts the bits left.
func (num *uint128le) shiftLeftNEW(p uint) {
	num << p
}

// Method shiftRight shifts the bits right.
func (num *uint128le) shiftRightNEW(p uint) {
	num >> p
}

// Method String returns 'num' as a string using shifts and a mask
func (num uint128le) String() string {
	s := shifts[2]
	var a [128]byte
	i := 0
	m := uintptr(1)
	for j := 0; j < len(num); j++ {
		for l := 0; l < 8; l++ {
			for num[j] > 1 {
				a[i] = digits[uintptr(num[j])&m]
				num[j] >>= s
				i++
				l++
			}
			if l < 7 && num[j] == 1 {
				a[i] = digits[uintptr(num[j])]
				num[j]--
			} else {
				a[i] = digits[uintptr(num[j])]
			}
			i++
		}
	}
	result := string(a[0:])
	return result
}

func (bstr *Bitstring) pad(l int) {
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = 48
	}
	*bstr = Bitstring(b) + *bstr
}

func fromBitstring(b Bitstring) uint128le {
	if len(b) < 128 {
		b.pad(128 - len(b))
	}
	var result uint128le
	var j int = 0
	sl := make([]byte, 8)
	var tempNum byte = 0
	pt := &tempNum
	for i := 0; i < 128; {
		for p, t := range b[i : i+8] {
			if t == 49 {
				sl[p] = 1
			} else {
				sl[p] = 0
			}
		}
		*pt = 0
		for q := 7; q >= 0; q-- {
			if q == 7 && sl[q] == 0 {
				continue
			} else if sl[q] == 1 {
				*pt += 1 << uint(q)
			}
		}
		result[j] = tempNum
		j++
		i += 8
	}

	return result
}
