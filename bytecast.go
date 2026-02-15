package bytecast

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"reflect"
)

// ToTypedValue converts reflect.Value type to it's underlying (real) type
//
// This function is useful when dealing with reflected value and converting it to byte array, sequence of actions in this case:
//
//	reflect.Value -> ToTypedValue() -> ToXBytes() -> [X]byte array
//
// Usage example:
//
//	https://stackoverflow.com/questions/78390431/simplify-making-a-generic-typed-value-from-a-reflect-value
func ToTypedValue[T any](v reflect.Value) (T, error) {
	u, ok := v.Interface().(T)
	if !ok {
		var t T // just a zero value
		return t, fmt.Errorf("v has not type %T", t)
	}
	return u, nil
}

// IntXXToBytesAndExpandWidth
//
//	Takes value passed in int64 container, but handles it as xx-bit value, represents as bytes slice and expand to specified width.
//
//	NOTE:
//	xx => bits
//	width => bytes
func IntXXToBytesAndExpandWidth(value int64, xx int, width int) ([]byte, error) {
	if xx <= 0 || xx > 64 {
		return nil, fmt.Errorf("unsupported bit size %d, must be 1..64", xx)
	}

	// [ s xxxxxxx xxxxxxxx ... xxxxxxxx ]   ← 56 bits
	//   ^
	//   sign bit (bit 55)

	maxPossibleV := int64(1)<<(xx-1) - 1
	minPossibleV := int64(-1) << (xx - 1)

	if value < minPossibleV || value > maxPossibleV {
		return nil, fmt.Errorf("value %d does not fit in int%d", value, xx)
	}

	// take lower 56 bits
	// (1 << 56)     = 100000...000
	// (1 << 56) - 1 = 011111...111
	var u uint64
	if value < 0 {
		u = uint64(value + (1 << xx)) // two's complement: 2^N + v (де v < 0)
	} else {
		u = uint64(value)
	}

	// How many bytes "netto" we need to store value?
	neededBytesNum := int(math.Ceil(float64(xx) / float64(8)))

	out := make([]byte, width)

	// Write into last neededBytesNum bytes (big-endian)
	// Example for int56 (neededBytesNum == 7) and width == 32:
	// 	out[25] = byte(u >> 48)
	//	out[26] = byte(u >> 40)
	//	out[27] = byte(u >> 32)
	//	out[28] = byte(u >> 24)
	//	out[29] = byte(u >> 16)
	//	out[30] = byte(u >> 8)
	//	out[31] = byte(u)
	offset := (neededBytesNum - 1) * 8
	for i := width - neededBytesNum; i < width; i++ {
		out[i] = byte(u >> uint(offset))
		offset -= 8
	}

	// sign extension (only for negative values)
	if value < 0 {
		for i := 0; i < width-neededBytesNum; i++ {
			out[i] = byte(0xff)
		}
	}

	return out, nil
}

// UintXXToBytesAndExpandWidth
//
//	Takes value passed in uint64 container, but handles it as xx-bit unsigned value,
//	represents it as bytes slice and expands to specified width (e.g., 32 bytes for EVM).
//
//	NOTE:
//	xx => bits
//	width => bytes
func UintXXToBytesAndExpandWidth(value uint64, xx int, width int) ([]byte, error) {
	if xx <= 0 || xx > 64 {
		return nil, fmt.Errorf("unsupported bit size %d, must be 1..64", xx)
	}

	// Перевіряємо, чи число поміщається в xx біт
	maxPossibleV := uint64(1)<<xx - 1
	if value > maxPossibleV {
		return nil, fmt.Errorf("value %d does not fit in uint%d", value, xx)
	}

	// Відкидаємо старші біти, залишаємо тільки xx бітів
	u := value & maxPossibleV

	// Скільки байт реально потрібно для зберігання xx бітів?
	neededBytesNum := int(math.Ceil(float64(xx) / 8.0))

	out := make([]byte, width)

	// Записуємо в останні neededBytesNum байт (big-endian)
	offset := (neededBytesNum - 1) * 8
	for i := width - neededBytesNum; i < width; i++ {
		out[i] = byte(u >> uint(offset))
		offset -= 8
	}

	// Для unsigned чисел старші байти просто нулі (у нас out вже zeroed)
	return out, nil
}

// IntXXFromBytes
//
//	Takes "bytes" bytes from input and interpret them as int-xx value,
//	then pack into int64 container and return.
func IntXXFromBytes(bytes []byte, xx int) (int64, error) {
	if xx <= 0 || xx > 64 {
		return 0, fmt.Errorf("unsupported bit size %d, must be 1..64", xx)
	}

	// Скільки байт реально потрібно для зберігання xx бітів?
	neededBytesNum := int(math.Ceil(float64(xx) / 8.0))

	if len(bytes) < neededBytesNum {
		return 0, fmt.Errorf(
			"expected at least %d bytes to interpret as int%d value, but got only %d bytes",
			neededBytesNum, xx, len(bytes),
		)
	}

	// We take last "neededBytes" bytes from received "bytes" bytes:
	var u uint64
	offset := (neededBytesNum - 1) * 8
	start := len(bytes) - neededBytesNum

	for i := start; i < len(bytes); i++ {
		u |= uint64(bytes[i]) << offset
		offset -= 8
	}

	signBitPosition := xx - 1

	// For int64 2s complement not needed, int64 already has correct sign
	if xx == 64 || u&(1<<signBitPosition) == 0 {
		return int64(u), nil
	}

	// EXAMPLE for int24:
	// if sign bit #24 is set → negative number
	//
	// 00000000 sxxxxxxx xxxxxxxx xxxxxxxx
	// ↑        ↑
	// біт 31   біт 23 (sign bit int24)
	//
	// v before:
	// 00000000 1xxxxxxx xxxxxxxx xxxxxxxx
	//
	// mask:
	// 11111111 00000000 00000000 00000000
	// ----------------------------------
	// v after OR:
	// 11111111 1xxxxxxx xxxxxxxx xxxxxxxx

	var mask uint64
	mask = ^((1 << xx) - 1)
	u |= mask

	return int64(u), nil
}

// UintXXFromBytes
//
//	Takes "bytes" bytes from input and interpret them as uint-xx value,
//	then pack into uint64 container and return.
func UintXXFromBytes(bytes []byte, xx int) (uint64, error) {
	if xx <= 0 || xx > 64 {
		return 0, fmt.Errorf("unsupported bit size %d, must be 1..64", xx)
	}

	// How many bytes needed to store xx bits?
	neededBytesNum := int(math.Ceil(float64(xx) / 8.0))

	if len(bytes) < neededBytesNum {
		return 0, fmt.Errorf(
			"expected at least %d bytes to interpret as uint%d value, but got only %d bytes",
			neededBytesNum, xx, len(bytes),
		)
	}

	// Take last neededBytesNum bytes (big-endian)
	var u uint64
	offset := (neededBytesNum - 1) * 8
	start := len(bytes) - neededBytesNum

	for i := start; i < len(bytes); i++ {
		u |= uint64(bytes[i]) << offset
		offset -= 8
	}

	return u, nil
}

// Int64To8Bytes
//
//	https://groups.google.com/g/golang-nuts/c/q1wk1WDNoo4?pli=1
func Int64To8Bytes(intValue int64) [8]byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(intValue))
	bFixed := (*[8]byte)(b)
	return *bFixed
}

func Int64ToBytesAndExpandWidth(intValue int64, width int) ([]byte, error) {
	if width < 8 {
		return []byte{}, fmt.Errorf("failed to convert int64 to bytes, provided width too short, got %d expected min 8", width)
	}

	byteValue := Int64To8Bytes(intValue)

	if intValue >= 0 {
		return LeftPadBytes00(byteValue[:], width), nil
	}

	return LeftPadBytesFF(byteValue[:], width), nil
}

func Int64From8Bytes(byteValue [8]byte) int64 {
	bSlice := byteValue[:]
	v := int64(binary.BigEndian.Uint64(bSlice))
	return v
}

func Int32To4Bytes(intValue int32) [4]byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(intValue))
	bFixed := (*[4]byte)(b)
	return *bFixed
}

func Int32ToBytesAndExpandWidth(intValue int32, width int) ([]byte, error) {
	if width < 4 {
		return []byte{}, fmt.Errorf("failed to convert int32 to bytes, provided width too short, got %d expected min 4", width)
	}

	byteValue := Int32To4Bytes(intValue)

	if intValue >= 0 {
		return LeftPadBytes00(byteValue[:], width), nil
	}

	return LeftPadBytesFF(byteValue[:], width), nil
}

func Int32From4Bytes(byteValue [4]byte) int32 {
	bSlice := byteValue[:]
	v := int32(binary.BigEndian.Uint32(bSlice))
	return v
}

func Uint32To4Bytes(intValue uint32) [4]byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, intValue)
	bFixed := (*[4]byte)(b)
	return *bFixed
}

func Uint32ToBytesAndExpandWidth(intValue uint32, width int) ([]byte, error) {
	if width < 4 {
		return []byte{}, fmt.Errorf("failed to convert uint32 to bytes, provided width too short, got %d expected min 4", width)
	}

	byteValue := Uint32To4Bytes(intValue)

	return LeftPadBytes00(byteValue[:], width), nil
}

func Uint32From4Bytes(byteValue [4]byte) uint32 {
	bSlice := byteValue[:]
	v := binary.BigEndian.Uint32(bSlice)
	return v
}

func Int16To2Bytes(intValue int16) [2]byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(intValue))
	bFixed := (*[2]byte)(b)
	return *bFixed
}

func Int16ToBytesAndExpandWidth(intValue int16, width int) ([]byte, error) {
	if width < 2 {
		return []byte{}, fmt.Errorf("failed to convert int16 to bytes, provided width too short, got %d expected min 2", width)
	}

	byteValue := Int16To2Bytes(intValue)

	if intValue >= 0 {
		return LeftPadBytes00(byteValue[:], width), nil
	}

	return LeftPadBytesFF(byteValue[:], width), nil
}

func Int16From2Bytes(byteValue [2]byte) int16 {
	bSlice := byteValue[:]
	v := int16(binary.BigEndian.Uint16(bSlice))
	return v
}

func Uint16To2Bytes(intValue uint16) [2]byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, intValue)
	bFixed := (*[2]byte)(b)
	return *bFixed
}

func Uint16ToBytesAndExpandWidth(intValue uint16, width int) ([]byte, error) {
	if width < 2 {
		return []byte{}, fmt.Errorf("failed to convert uint16 to bytes, provided width too short, got %d expected min 2", width)
	}

	byteValue := Uint16To2Bytes(intValue)

	return LeftPadBytes00(byteValue[:], width), nil
}

func Uint16From2Bytes(byteValue [2]byte) uint16 {
	bSlice := byteValue[:]
	v := binary.BigEndian.Uint16(bSlice)
	return v
}

func Int8To1Byte(intValue int8) [1]byte {
	return [1]byte{byte(intValue)}
}

func Int8ToBytesAndExpandWidth(intValue int8, width int) ([]byte, error) {
	if width < 1 {
		return []byte{}, fmt.Errorf("failed to convert int8 to bytes, provided width too short, got %d expected min 1", width)
	}

	byteValue := Int8To1Byte(intValue)

	if intValue >= 0 {
		return LeftPadBytes00(byteValue[:], width), nil
	}

	return LeftPadBytesFF(byteValue[:], width), nil
}

func Int8From1Byte(byteValue [1]byte) int8 {
	return int8(byteValue[0])
}

func Uint8To1Byte(intValue uint8) [1]byte {
	return [1]byte{intValue}
}

func Uint8ToBytesAndExpandWidth(intValue uint8, width int) ([]byte, error) {
	if width < 1 {
		return []byte{}, fmt.Errorf("failed to convert uint8 to bytes, provided width too short, got %d expected min 1", width)
	}

	byteValue := Uint8To1Byte(intValue)

	return LeftPadBytes00(byteValue[:], width), nil
}

func Uint8From1Byte(byteValue [1]byte) uint8 {
	return byteValue[0]
}

func BigIntToBytesAndExpandWidth(bigInt *big.Int, width int) ([]byte, error) {
	if bigInt == nil {
		bigInt = big.NewInt(0)
	}

	bitLen := uint(width * 8)

	if bigInt.Sign() >= 0 {
		if bigInt.BitLen() > int(bitLen) {
			return nil, fmt.Errorf("integer %s too large to encode in %d bytes", bigInt, width)
		}
		return LeftPadBytes00(bigInt.Bytes(), width), nil
	}

	// two's complement: 2^N + v (де v < 0)
	// Example:
	// 	-1 = 2^256 - 1
	//	-2 = 2^256 - 2
	mod := new(big.Int).Lsh(big.NewInt(1), bitLen)
	twos := new(big.Int).Add(mod, bigInt)

	if twos.Sign() < 0 || twos.BitLen() > int(bitLen) {
		return nil, fmt.Errorf("integer %s cannot fit in %d bytes", bigInt.String(), width)
	}

	// вже правильне представлення у width байт, не добиваємо нулями
	return twos.Bytes(), nil
}

func BigIntFromBytes(byteValue []byte) *big.Int {
	// unsigned interpretation
	x := new(big.Int).SetBytes(byteValue)

	// if sign bit is NOT set → positive
	if len(byteValue) == 0 || byteValue[0]&0x80 == 0 { // 0x80 is highest, sign bit, the same as 0b10000000
		return x
	}

	// negative number:
	// x = x - 2^(8 * len(b))
	bitLen := uint(len(byteValue) * 8)
	mod := big.NewInt(0).Lsh(big.NewInt(1), bitLen)

	return x.Sub(x, mod)
}

// BigIntXXXFromBytes
//
//	Takes "bytes" bytes from input and interpret them as big.Int-xxx (int128, int256) value,
//	then pack into *big.Int and return.
//
//	This method is better than BigIntFromBytes,
//	because it CHECKS SIGN BIT AT STRICTLY DEFINED PLACE.
func BigIntXXXFromBytes(bytes []byte, xxx int) (*big.Int, error) {
	if xxx <= 64 {
		return nil, fmt.Errorf("too small bit size %d, for standard int up to int64 use \"IntXXFromBytes\" method", xxx)
	}

	// Скільки байт реально потрібно для зберігання xx бітів?
	neededBytesNum := int(math.Ceil(float64(xxx) / 8.0))

	if len(bytes) < neededBytesNum {
		return nil, fmt.Errorf(
			"expected at least %d bytes to interpret as big.Int (%dbit) value, but got only %d bytes",
			neededBytesNum, xxx, len(bytes),
		)
	}

	// We take last "neededBytes" bytes from received "bytes" bytes:
	start := len(bytes) - neededBytesNum
	uBytes := bytes[start:]

	resultUnsigned := new(big.Int).SetBytes(uBytes)

	// signBitPosition = xxx - 1; if sign bit is NOT set → positive
	if resultUnsigned.Bit(xxx-1) == 0 {
		return resultUnsigned, nil
	}

	// negative number:
	// x = x - 2^xxx
	mod := new(big.Int).Lsh(big.NewInt(1), uint(xxx))
	return resultUnsigned.Sub(resultUnsigned, mod), nil
}

func BoolTo1Byte(boolVal bool) [1]byte {
	if boolVal {
		return [1]byte{0x01} // завжди 1 для true
	}
	return [1]byte{0x00} // завжди 0 для false
}

func BoolToBytesAndExpandWidth(boolVal bool, width int) ([]byte, error) {
	if width < 1 {
		return []byte{}, fmt.Errorf("failed to convert bool to bytes, provided width too short, got %d expected min 1", width)
	}

	byteValue := BoolTo1Byte(boolVal)

	return LeftPadBytes00(byteValue[:], width), nil
}

func BoolFrom1Byte(bytesVal [1]byte) bool {
	return bytesVal[0] != 0 // будь-яке ненульове значення → true
}

// StringTo256Bytes converts arbitrary string to bytes array.
//
//	IMPORTANT! Max input string length LIMITED TO 255 bytes
//	(this is 255 ascii symbols where 1 symbol can be represented by 1 byte)
func StringTo256Bytes(stringValue string) ([256]byte, error) {
	stringBytes := []byte(stringValue)
	l := len(stringBytes)

	if l > 255 {
		return [256]byte{}, fmt.Errorf("string length exceeded, max 255 bytes allowed")
	}

	fixedWidthData := make([]byte, 0, 256)
	fixedWidthData = append(fixedWidthData, uint8(l))
	fixedWidthData = append(fixedWidthData, LeftPadBytes00(stringBytes, 255)...)

	dataArray := (*[256]byte)(fixedWidthData)

	return *dataArray, nil
}

func StringFrom256Bytes(byteVal [256]byte) string {
	significantBytesCount := uint8(byteVal[0])

	if significantBytesCount == 0 {
		return ""
	}

	significantBytes := byteVal[256-int(significantBytesCount):]

	return string(significantBytes)
}

func LeftPadBytes(input []byte, size int, padByte byte) []byte {
	if size <= len(input) {
		return input
	}
	pad := bytes.Repeat([]byte{padByte}, size-len(input))
	return append(pad, input...)
}

func LeftPadBytes00(input []byte, size int) []byte {
	return LeftPadBytes(input, size, 0x00)
}

func LeftPadBytesFF(input []byte, size int) []byte {
	return LeftPadBytes(input, size, 0xFF)
}
