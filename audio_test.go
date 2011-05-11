package captcha

import (
	"io/ioutil"
	"testing"
)

func BenchmarkNewAudio(b *testing.B) {
	b.StopTimer()
	d := RandomDigits(DefaultLen)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		NewAudio(d, "")
	}
}

func BenchmarkAudioWriteTo(b *testing.B) {
	b.StopTimer()
	d := RandomDigits(DefaultLen)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a := NewAudio(d, "")
		n, _ := a.WriteTo(ioutil.Discard)
		b.SetBytes(n)
	}
}
