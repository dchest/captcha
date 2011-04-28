package captcha

import (
	"os"
	"testing"
)

type byteCounter struct {
	n int64
}

func (bc *byteCounter) Write(b []byte) (int, os.Error) {
	bc.n += int64(len(b))
	return len(b), nil
}

func BenchmarkNewImage(b *testing.B) {
	b.StopTimer()
	d := RandomDigits(DefaultLen)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		NewImage(d, StdWidth, StdHeight)
	}
}

func BenchmarkImageWriteTo(b *testing.B) {
	b.StopTimer()
	d := RandomDigits(DefaultLen)
	b.StartTimer()
	counter := &byteCounter{}
	for i := 0; i < b.N; i++ {
		img := NewImage(d, StdWidth, StdHeight)
		img.WriteTo(counter)
		b.SetBytes(counter.n)
		counter.n = 0
	}
}
