package captcha

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"math"
	"os"
	"rand"
	"io"
)

const sampleRate = 8000

var (
	// Length of the longest number sound
	longestNumSndLen int
	endingBeepSound  []byte
)

func init() {
	for _, v := range numberSounds {
		if longestNumSndLen < len(v) {
			longestNumSndLen = len(v)
		}
	}
	endingBeepSound = changeSpeed(beepSound, 1.4)
}

type Audio struct {
	body *bytes.Buffer
}

// NewImage returns a new audio captcha with the given slice of numbers, where
// each number must be in range 0-9.
func NewAudio(numbers []byte) *Audio {
	numsnd := make([][]byte, len(numbers))
	nsdur := 0
	for i, n := range numbers {
		snd := randomizedNumSound(n)
		nsdur += len(snd)
		numsnd[i] = snd
	}
	// Intervals between numbers (including beginning)
	intervals := make([]int, len(numbers)+1)
	intdur := 0
	for i := range intervals {
		// 1 to 3 seconds
		dur := rnd(sampleRate, sampleRate*3)
		intdur += dur
		intervals[i] = dur
	}
	// Background noise
	bg := makeBackgroundSound(longestNumSndLen*len(numbers) + intdur)
	// --
	a := new(Audio)
	a.body = bytes.NewBuffer(nil)
	// Prelude, three beeps
	sil := makeSilence(sampleRate / 5)
	a.body.Write(beepSound)
	a.body.Write(sil)
	a.body.Write(beepSound)
	a.body.Write(sil)
	a.body.Write(beepSound)
	// Numbers
	pos := intervals[0]
	for i, v := range numsnd {
		mixSound(bg[pos:], v)
		pos += len(v) + intervals[i+1]
	}
	a.body.Write(bg)
	// Ending
	a.body.Write(endingBeepSound)
	return a
}

// NewRandomAudio generates a sequence of random numbers with the given length,
// and returns a new audio captcha with this numbers, and the sequence of
// numbers itself.
func NewRandomAudio(length int) (a *Audio, numbers []byte) {
	numbers = randomNumbers(length)
	a = NewAudio(numbers)
	return
}

// WriteTo writes captcha audio in WAVE format into the given io.Writer, and
// returns the number of bytes written and an error if any.
func (a *Audio) WriteTo(w io.Writer) (n int64, err os.Error) {
	nn, err := w.Write(waveHeader)
	n = int64(nn)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.LittleEndian, uint32(a.body.Len()))
	if err != nil {
		return
	}
	nn += 4
	n, err = a.body.WriteTo(w)
	n += int64(nn)
	return
}

// mixSound mixes src into dst. Dst must have length equal to or greater than
// src length.
func mixSound(dst, src []byte) {
	for i, v := range src {
		av := int(v)
		bv := int(dst[i])
		if av < 128 && bv < 128 {
			dst[i] = byte(av * bv / 128)
		} else {
			dst[i] = byte(2*(av+bv) - av*bv/128 - 256)
		}
	}
}

func setSoundLevel(a []byte, level float64) {
	for i, v := range a {
		av := float64(v)
		switch {
		case av > 128:
			if av = (av-128)*level + 128; av < 128 {
				av = 128
			}
		case av < 128:
			if av = 128 - (128-av)*level; av > 128 {
				av = 128
			}
		default:
			continue
		}
		a[i] = byte(av)
	}
}

// changeSpeed returns new PCM bytes from the bytes with the speed and pitch
// changed to the given value that must be in range [0, x].
func changeSpeed(a []byte, pitch float64) []byte {
	b := make([]byte, int(math.Floor(float64(len(a))*pitch)))
	var p float64
	for _, v := range a {
		for i := int(p); i < int(p+pitch); i++ {
			b[i] = v
		}
		p += pitch
	}
	return b
}

// rndFloat64n returns a random float64 number in range [from, to].
func rndFloat64n(from, to float64) float64 {
	return (to-from)*rand.Float64() + from
}

func randomSpeed(a []byte) []byte {
	pitch := rndFloat64n(0.9, 1.2)
	return changeSpeed(a, pitch)
}

func makeSilence(length int) []byte {
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = 128
	}
	return b
}

func makeStaticNoise(length int, level uint8) []byte {
	noise := make([]byte, length)
	_, err := io.ReadFull(crand.Reader, noise)
	if err != nil {
		panic("error reading from random source: " + err.String())
	}
	for i := 0; i < len(noise); i++ {
		noise[i] %= level
		noise[i] += 128 - level/2
	}
	return noise
}

func reversedSound(a []byte) []byte {
	ln := len(a)
	b := make([]byte, ln)
	for i, v := range a {
		b[ln-1-i] = v
	}
	return b
}

func makeBackgroundSound(length int) []byte {
	b := makeStaticNoise(length, 8)
	for i := 0; i < length/(sampleRate/10); i++ {
		snd := numberSounds[rand.Intn(10)]
		snd = changeSpeed(reversedSound(snd), rndFloat64n(0.8, 1.4))
		place := rand.Intn(len(b) - len(snd))
		setSoundLevel(snd, rndFloat64n(0.5, 1.2))
		mixSound(b[place:], snd)
	}
	setSoundLevel(b, rndFloat64n(0.2, 0.3))
	return b
}

func randomizedNumSound(n byte) []byte {
	s := randomSpeed(numberSounds[n])
	setSoundLevel(s, rndFloat64n(0.7, 1.3))
	return s
}
