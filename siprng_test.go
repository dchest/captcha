package captcha

import (
	"bytes"
	"testing"
)

func TestSiphash(t *testing.T) {
	good := uint64(0xe849e8bb6ffe2567)
	cur := siphash(0, 0, 0)
	if cur != good {
		t.Fatalf("siphash: expected %x, got %x", good, cur)
	}
}

func TestSiprng(t *testing.T) {
	m := make(map[uint64]interface{})
	var yes interface{}
	r := siprng{}
	r.Seed([16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	for i := 0; i < 100000; i++ {
		v := r.Uint64()
		if _, ok := m[v]; ok {
			t.Errorf("siphash: collision on %d: %x", i, v)
		}
		m[v] = yes
	}
}

func TestSiprngBytes(t *testing.T) {
	r := siprng{}
	r.Seed([16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	x := r.Bytes(32)
	if len(x) != 32 {
		t.Fatalf("siphash: wrong length: expected 32, got %d", len(x))
	}
	y := r.Bytes(32)
	if bytes.Equal(x, y) {
		t.Fatalf("siphash: stream repeats: %x = %x", x, y)
	}
	r.Seed([16]byte{})
	z := r.Bytes(32)
	if bytes.Equal(z, x) {
		t.Fatalf("siphash: outputs under different keys repeat: %x = %x", z, x)
	}
}

func BenchmarkSiprng(b *testing.B) {
	b.SetBytes(8)
	p := &siprng{}
	for i := 0; i < b.N; i++ {
		p.Uint64()
	}
}
