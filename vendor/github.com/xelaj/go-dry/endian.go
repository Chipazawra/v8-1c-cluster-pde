package dry

import (
	"math"
	"unsafe"
)

func EndianIsLittle() bool {
	var word uint16 = 1
	littlePtr := (*uint8)(unsafe.Pointer(&word))
	return (*littlePtr) == 1
}

func EndianIsBig() bool {
	return !EndianIsLittle()
}

func EndianSafeSplitUint16(value uint16) (leastSignificant, mostSignificant uint8) {
	bytes := (*[2]uint8)(unsafe.Pointer(&value))
	if EndianIsLittle() {
		return bytes[0], bytes[1]
	}
	return bytes[1], bytes[0]
}

// TODO: короче, проблема в представлении интов как последовательности байт, и хер поймет как правильно отрицательные числа переводить в негативные
// https://www.rapidtables.com/convert/number/decimal-to-hex.html
func NormalizeInt(n int) uint32 {
	return uint32(0x100000000 + n)
}

// не работает // FIXME
func NormalizeLong(n int64) uint64 {
	return uint64(math.MaxInt64 + n)
}
