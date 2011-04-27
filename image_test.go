package captcha

import (
	"os"
	"testing"
)

type devNull struct{}

func (devNull) Write(b []byte) (int, os.Error) {
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
	for i := 0; i < b.N; i++ {
		img := NewImage(d, StdWidth, StdHeight)
		img.WriteTo(devNull{}) //TODO(dchest): use ioutil.Discard when its available
	}
}
