package bytecast

import (
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

func Uint16From2Bytes(byteValue [2]byte) uint16 {
	bSlice := byteValue[:]
	v := binary.BigEndian.Uint16(bSlice)
	return v
}

func Int8To1Byte(intValue int8) [1]byte {
	return [1]byte{byte(intValue)}
}

func Int8From1Byte(byteValue [1]byte) int8 {
	return int8(byteValue[0])
}

func Uint8To1Byte(intValue uint8) [1]byte {
	return [1]byte{intValue}
}

func Uint8From1Byte(byteValue [1]byte) uint8 {
	return byteValue[0]
}

func BigIntTo32Bytes(bigInt *big.Int) [32]byte {
	if bigInt == nil {
		bigInt = big.NewInt(0)
	}

	bigInt32Bytes := LeftPadBytes(bigInt.Bytes(), 32)
	bigInt32BFixedArray := ([32]byte)(bigInt32Bytes) // Slice to array (array pointer) conversion

	return bigInt32BFixedArray
}

func BigIntFrom32Bytes(byteValue [32]byte) *big.Int {
	bigInt := big.NewInt(0).SetBytes(byteValue[:])

	return bigInt
}

func BoolTo1Byte(boolVal bool) [1]byte {
	valueInt8 := int8(0)
	if boolVal {
		valueInt8 = 1
	}

	bytesArray := [1]byte{byte(valueInt8)}

	return bytesArray
}

func BoolFrom1Byte(bytesVal [1]byte) bool {
	var valueInt8 int8
	valueInt8 = int8(bytesVal[0])

	if valueInt8 > 0 {
		return true
	}

	return false
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
	fixedWidthData = append(fixedWidthData, LeftPadBytes(stringBytes, 255)...)

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

func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}
