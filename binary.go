package serpent

type uint128le [16]byte

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

func newFromUint64(num uint64) (result uint128le) {
	for i := 0; i < 8; i++ {
		result[i] = byte(num & 0xff)
		num >>= 8
	}
	return result
}

func (num *uint128le) binaryXor(other uint128le) {
	for i := 0; i < 16; i++ {
		num[i] ^= other[i]
	}
}

func xor(args []uint128le) *uint128le {
	result := &args[0]
	for i := 1; i < len(args); i++ {
		result.binaryXor(args[i])
	}
	return result
}

// Return 'num' as a string using shifts
func (num uint128le) String() string {
	b := uint64(2)
	s := shifts[2]
	var a [128 + 1]byte
	i := len(a)
	m := uintptr(b) - 1
	for j := 0; j < len(num); j++ {
		for uint64(num[j]) >= b {
			i--
			a[i] = digits[uintptr(num[j])&m]
			num[j] >>= s
		}
		i--
		a[i] = digits[uintptr(num[j])]
	}
	result := string(a[i:])
	return result
}
