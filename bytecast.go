package bytecast

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
