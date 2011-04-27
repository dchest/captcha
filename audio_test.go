package captcha

import "testing"

func BenchmarkNewAudio(b *testing.B) {
	b.StopTimer()
	d := RandomDigits(DefaultLen)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		NewAudio(d)
	}
}

func BenchmarkAudioWriteTo(b *testing.B) {
	b.StopTimer()
	d := RandomDigits(DefaultLen)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a := NewAudio(d)
		a.WriteTo(devNull{}) //TODO(dchest): use ioutil.Discard when its available
	}
}
