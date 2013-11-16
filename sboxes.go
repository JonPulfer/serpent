package serpent

import ()

// S0: 3  8 15  1 10  6  5 11 14 13  4  2  7  0  9 12
// depth = 5,7,4,2, Total gates=18
func RND00(a, b, c, d, w, x, y, z *uint32le) {
	var t02, t03, t05, t06, t07, t08, t09, t11, t12, t13, t14, t15, t17, t01 uint32le
	t01 = b   ^ c
	t02 = a   | d
	t03 = a   ^ b
	z   = t02 ^ t01
	t05 = c   | z
	t06 = a   ^ d
	t07 = b   | c
	t08 = d   & t05
	t09 = t03 & t07
	y   = t09 ^ t08
	t11 = t09 & y
	t12 = c   ^ d
	t13 = t07 ^ t11
	t14 = b   & t06
	t15 = t06 ^ t13
	w   =     ~ t15
	t17 = w   ^ t14
	x   = t12 ^ t17
}

func InvRND00(a, b, c, d, w, x, y, z *uint32le) {
	var t02, t03, t04, t05, t06, t08, t09, t10, t12, t13, t14, t15, t17, t18, t01 uint32le
	t01 = c   ^ d
	t02 = a   | b
	t03 = b   | c
	t04 = c   & t01
	t05 = t02 ^ t01
	t06 = a   | t04
	y   =     ~ t05
	t08 = b   ^ d
	t09 = t03 & t08
	t10 = d   | y
	x   = t09 ^ t06
	t12 = a   | t05
	t13 = x   ^ t12
	t14 = t03 ^ t10
	t15 = a   ^ c
	z   = t14 ^ t13
	t17 = t05 & t13
	t18 = t14 | t17
	w   = t15 ^ t18
}
