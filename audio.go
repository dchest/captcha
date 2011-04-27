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

const sampleRate = 8000 // Hz

var (
	longestDigitSndLen int
	endingBeepSound    []byte
	reverseDigitSounds [][]byte
)

func init() {
	for _, v := range digitSounds {
		if longestDigitSndLen < len(v) {
			longestDigitSndLen = len(v)
		}
	}
	endingBeepSound = changeSpeed(beepSound, 1.4)
	// Preallocate reversed digit sounds for background noise.
	reverseDigitSounds = make([][]byte, len(digitSounds))
	for i, v := range digitSounds {
		reverseDigitSounds[i] = reversedSound(v)
	}
}

// BUG(dchest): [Not our bug] Google Chrome 10 plays unsigned 8-bit PCM WAVE
// audio on Mac with horrible distortions.  Issue:
// http://code.google.com/p/chromium/issues/detail?id=70730.
// This has been fixed, and version 12 will play them properly.

type Audio struct {
	body *bytes.Buffer
}

// NewImage returns a new audio captcha with the given digits, where each digit
// must be in range 0-9.
func NewAudio(digits []byte) *Audio {
	numsnd := make([][]byte, len(digits))
	nsdur := 0
	for i, n := range digits {
		snd := randomizedDigitSound(n)
		nsdur += len(snd)
		numsnd[i] = snd
	}
	// Random intervals between digits (including beginning).
	intervals := make([]int, len(digits)+1)
	intdur := 0
	for i := range intervals {
		dur := rnd(sampleRate, sampleRate*3) // 1 to 3 seconds
		intdur += dur
		intervals[i] = dur
	}
	// Generate background sound.
	bg := makeBackgroundSound(longestDigitSndLen*len(digits) + intdur)
	// Create buffer and write audio to it.
	a := new(Audio)
	sil := makeSilence(sampleRate / 5)
	bufcap := 3*len(beepSound) + 2*len(sil) + len(bg) + len(endingBeepSound)
	a.body = bytes.NewBuffer(make([]byte, 0, bufcap))
	// Write prelude, three beeps.
	a.body.Write(beepSound)
	a.body.Write(sil)
	a.body.Write(beepSound)
	a.body.Write(sil)
	a.body.Write(beepSound)
	// Write digits.
	pos := intervals[0]
	for i, v := range numsnd {
		mixSound(bg[pos:], v)
		pos += len(v) + intervals[i+1]
	}
	a.body.Write(bg)
	// Write ending (one beep).
	a.body.Write(endingBeepSound)
	return a
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

// EncodedLen returns the length of WAV-encoded audio captcha.
func (a *Audio) EncodedLen() int {
	return len(waveHeader) + 4 + a.body.Len()
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
func changeSpeed(a []byte, speed float64) []byte {
	b := make([]byte, int(math.Floor(float64(len(a))*speed)))
	var p float64
	for _, v := range a {
		for i := int(p); i < int(p+speed); i++ {
			b[i] = v
		}
		p += speed
	}
	return b
}

// rndf returns a random float64 number in range [from, to].
func rndf(from, to float64) float64 {
	return (to-from)*rand.Float64() + from
}

func randomSpeed(a []byte) []byte {
	pitch := rndf(0.9, 1.2)
	return changeSpeed(a, pitch)
}

func makeSilence(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = 128
	}
	return b
}

func makeWhiteNoise(length int, level uint8) []byte {
	noise := make([]byte, length)
	if _, err := io.ReadFull(crand.Reader, noise); err != nil {
		panic("error reading from random source: " + err.String())
	}
	adj := 128 - level/2
	for i, v := range noise {
		v %= level
		v += adj
		noise[i] = v
	}
	return noise
}

func reversedSound(a []byte) []byte {
	n := len(a)
	b := make([]byte, n)
	for i, v := range a {
		b[n-1-i] = v
	}
	return b
}

func makeBackgroundSound(length int) []byte {
	b := makeWhiteNoise(length, 8)
	for i := 0; i < length/(sampleRate/10); i++ {
		snd := reverseDigitSounds[rand.Intn(10)]
		snd = changeSpeed(snd, rndf(0.8, 1.4))
		place := rand.Intn(len(b) - len(snd))
		setSoundLevel(snd, rndf(0.5, 1.2))
		mixSound(b[place:], snd)
	}
	setSoundLevel(b, rndf(0.2, 0.3))
	return b
}

func randomizedDigitSound(n byte) []byte {
	s := randomSpeed(digitSounds[n])
	setSoundLevel(s, rndf(0.7, 1.3))
	return s
}
