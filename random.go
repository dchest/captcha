// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	crand "crypto/rand"
	"io"
	"rand"
	"time"
)

// idLen is a length of captcha id string.
const idLen = 20

// idChars are characters allowed in captcha id.
var idChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func init() {
	rand.Seed(time.Nanoseconds())
}

// RandomDigits returns a byte slice of the given length containing
// pseudorandom numbers in range 0-9. The slice can be used as a captcha
// solution.
func RandomDigits(length int) (b []byte) {
	b = randomBytes(length)
	for i := range b {
		b[i] %= 10
	}
	return
}

// randomBytes returns a byte slice of the given length read from CSPRNG.
func randomBytes(length int) (b []byte) {
	b = make([]byte, length)
	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		panic("captcha: error reading random source: " + err.String())
	}
	return
}

// randomId returns a new random id string.
func randomId() string {
	b := randomBytes(idLen)
	alen := byte(len(idChars))
	for i, c := range b {
		b[i] = idChars[c%alen]
	}
	return string(b)
}

// rnd returns a non-crypto pseudorandom int in range [from, to].
func rnd(from, to int) int {
	return rand.Intn(to+1-from) + from
}

// rndf returns a non-crypto pseudorandom float64 in range [from, to].
func rndf(from, to float64) float64 {
	return (to-from)*rand.Float64() + from
}
