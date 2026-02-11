package bytecast

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestIntXXToBytesAndExpandWidth(t *testing.T) {
	tests := []struct {
		value int64
		bits  int
		want  string // hex representation
	}{
		{0, 24, "0000000000000000000000000000000000000000000000000000000000000000"},
		{1, 24, "0000000000000000000000000000000000000000000000000000000000000001"},
		{-1, 24, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{-193630, 24, "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd0ba2"},
		{8388607, 24, "00000000000000000000000000000000000000000000000000000000007fffff"},
		{-8388608, 24, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff800000"},

		{1234567890123, 56, "0000000000000000000000000000000000000000000000000000011f71fb04cb"},
		{-1234567890123, 56, "fffffffffffffffffffffffffffffffffffffffffffffffffffffee08e04fb35"},
		{36028797018963967, 56, "000000000000000000000000000000000000000000000000007fffffffffffff"},
		{-36028797018963968, 56, "ffffffffffffffffffffffffffffffffffffffffffffffffff80000000000000"},
	}

	for _, tt := range tests {
		out, err := IntXXToBytesAndExpandWidth(tt.value, tt.bits, 32)
		if err != nil {
			t.Errorf("IntXXToBytesAndExpandWidth(%d, %d) returned error: %v", tt.value, tt.bits, err)
			continue
		}

		got := fmt.Sprintf("%x", out)
		if got != tt.want {
			t.Errorf("IntXXToBytesAndExpandWidth(%d, %d) = %s; want %s", tt.value, tt.bits, got, tt.want)
		}
	}
}

func TestUintXXToBytesAndExpandWidth(t *testing.T) {
	tests := []struct {
		value uint64
		bits  int
		want  string // очікуване hex-представлення 32 байт
	}{
		// uint8
		{0, 8, "0000000000000000000000000000000000000000000000000000000000000000"},
		{255, 8, "00000000000000000000000000000000000000000000000000000000000000ff"},

		// uint16
		{0, 16, "0000000000000000000000000000000000000000000000000000000000000000"},
		{65535, 16, "000000000000000000000000000000000000000000000000000000000000ffff"},

		// uint24
		{0, 24, "0000000000000000000000000000000000000000000000000000000000000000"},
		{16777215, 24, "0000000000000000000000000000000000000000000000000000000000ffffff"},

		// uint32
		{0, 32, "0000000000000000000000000000000000000000000000000000000000000000"},
		{4294967295, 32, "00000000000000000000000000000000000000000000000000000000ffffffff"},

		// uint56
		{0, 56, "0000000000000000000000000000000000000000000000000000000000000000"},
		{72057594037927935, 56, "00000000000000000000000000000000000000000000000000ffffffffffffff"},

		// uint64
		{0, 64, "0000000000000000000000000000000000000000000000000000000000000000"},
		{18446744073709551615, 64, "000000000000000000000000000000000000000000000000ffffffffffffffff"},
	}

	for _, tt := range tests {
		out, err := UintXXToBytesAndExpandWidth(tt.value, tt.bits, 32)
		if err != nil {
			t.Errorf("UintXXToBytesAndExpandWidth(%d, %d) returned error: %v", tt.value, tt.bits, err)
			continue
		}

		got := fmt.Sprintf("%x", out)
		if got != tt.want {
			t.Errorf("UintXXToBytesAndExpandWidth(%d, %d) = %s; want %s", tt.value, tt.bits, got, tt.want)
		}
	}
}

func TestIntXXFromBytes(t *testing.T) {

	value_0_24, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	value_1_24, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	value_m1_24, _ := hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	value_m193630_24, _ := hex.DecodeString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd0ba2")
	value_8388607_24, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000007fffff")
	value_m8388608_24, _ := hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff800000")
	value_1234567890123_56, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000011f71fb04cb")
	value_m1234567890123_56, _ := hex.DecodeString("fffffffffffffffffffffffffffffffffffffffffffffffffffffee08e04fb35")
	value_0_56, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	value_36028797018963967_56, _ := hex.DecodeString("000000000000000000000000000000000000000000000000007fffffffffffff")
	value_m36028797018963968_56, _ := hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffff80000000000000")

	tests := []struct {
		bytes []byte
		bits  int
		want  int64
	}{
		{value_0_24, 24, 0},
		{value_1_24, 24, 1},

		{value_m1_24, 24, -1},
		{value_m193630_24, 24, -193630},

		{value_8388607_24, 24, 8388607},
		{value_m8388608_24, 24, -8388608},

		{value_1234567890123_56, 56, 1234567890123},
		{value_m1234567890123_56, 56, -1234567890123},

		{value_0_56, 56, 0},
		{value_36028797018963967_56, 56, 36028797018963967},

		{value_m36028797018963968_56, 56, -36028797018963968},
	}

	for _, tt := range tests {
		out, err := IntXXFromBytes(tt.bytes, tt.bits)
		if err != nil {
			t.Errorf("IntXXFromBytes(%x, %d) returned error: %v", tt.bytes, tt.bits, err)
			continue
		}

		if out != tt.want {
			t.Errorf("IntXXFromBytes(%x, %d) = %d; want %d", tt.bytes, tt.bits, out, tt.want)
		}
	}
}

func TestUintXXFromBytes(t *testing.T) {

	value_0_8, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	value_255_8, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000ff")

	value_0_16, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	value_65535_16, _ := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000ffff")

	value_0_24, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	value_16777215_24, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000ffffff")

	value_0_32, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	value_4294967295_32, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000ffffffff")

	value_0_56, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	value_72057594037927935_56, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000ffffffffffffff")

	value_0_64, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	value_18446744073709551615_64, _ := hex.DecodeString("000000000000000000000000000000000000000000000000ffffffffffffffff")

	tests := []struct {
		bytes []byte
		bits  int
		want  uint64
	}{
		// uint8
		{value_0_8, 8, 0},
		{value_255_8, 8, 255},

		// uint16
		{value_0_16, 16, 0},
		{value_65535_16, 16, 65535},

		// uint24
		{value_0_24, 24, 0},
		{value_16777215_24, 24, 16777215},

		// uint32
		{value_0_32, 32, 0},
		{value_4294967295_32, 32, 4294967295},

		// uint56
		{value_0_56, 56, 0},
		{value_72057594037927935_56, 56, 72057594037927935},

		// uint64
		{value_0_64, 64, 0},
		{value_18446744073709551615_64, 64, 18446744073709551615},
	}

	for _, tt := range tests {
		out, err := UintXXFromBytes(tt.bytes, tt.bits)
		if err != nil {
			t.Errorf("UintXXFromBytes(%x, %d) returned error: %v", tt.bytes, tt.bits, err)
			continue
		}

		if out != tt.want {
			t.Errorf("UintXXFromBytes(%x, %d) = %d; want %d", tt.bytes, tt.bits, out, tt.want)
		}
	}
}

func TestInt32SignExtension(t *testing.T) {
	cases := []int32{
		100,
		1,
		0,
		-1,
		-100,
		math.MaxInt32,
		math.MinInt32,
	}

	for _, v := range cases {
		b, err := Int32ToBytesAndExpandWidth(v, 32)
		if err != nil {
			t.Fatal(err)
		}

		got := BigIntFromBytes(b)
		if got.Int64() != int64(v) {
			t.Fatalf("expected %d got %d", v, got.Int64())
		}
	}
}

func TestBigIntNegativeRoundTrip(t *testing.T) {
	cases := []*big.Int{
		big.NewInt(-1),
		big.NewInt(-2),
		big.NewInt(-255),
		big.NewInt(-256),
		new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 255)),
	}

	for _, v := range cases {
		b, err := BigIntToBytesAndExpandWidth(v, 32)
		if err != nil {
			t.Fatal(err)
		}

		got := BigIntFromBytes(b)
		if got.Cmp(v) != 0 {
			t.Fatalf("expected %s got %s", v, got)
		}
	}
}

func TestWidthTooSmall(t *testing.T) {
	_, err := Int32ToBytesAndExpandWidth(1, 3)
	if err == nil {
		t.Fatal("expected error")
	}

	_, err = BigIntToBytesAndExpandWidth(big.NewInt(256), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBoolFromByte(t *testing.T) {
	cases := map[byte]bool{
		0x00: false,
		0x01: true,
		0x02: true,
		0xFF: true,
	}

	for b, expected := range cases {
		got := BoolFrom1Byte([1]byte{b})
		if got != expected {
			t.Fatalf("byte %x expected %v got %v", b, expected, got)
		}
	}
}

func TestBoolFromByteFullRange(t *testing.T) {
	for b := 0; b <= 0xFF; b++ {
		got := BoolFrom1Byte([1]byte{byte(b)})
		expected := b != 0 // будь-яке ненульове значення → true

		if got != expected {
			t.Fatalf("byte 0x%02X: expected %v, got %v", b, expected, got)
		}
	}
}

func TestBoolRoundTrip(t *testing.T) {
	cases := []bool{true, false}

	for _, val := range cases {
		b := BoolTo1Byte(val)
		got := BoolFrom1Byte(b)

		if got != val {
			t.Fatalf("round-trip failed for %v: got %v", val, got)
		}
	}
}

func TestString256Boundaries(t *testing.T) {
	s := strings.Repeat("a", 255)
	b, err := StringTo256Bytes(s)
	if err != nil {
		t.Fatal(err)
	}

	got := StringFrom256Bytes(b)
	if got != s {
		t.Fatalf("mismatch")
	}

	_, err = StringTo256Bytes(strings.Repeat("a", 256))
	if err == nil {
		t.Fatal("expected error")
	}
}
