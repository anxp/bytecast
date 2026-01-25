package bytecast

import (
	"math"
	"math/big"
	"strings"
	"testing"
)

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
