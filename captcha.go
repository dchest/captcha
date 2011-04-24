package captcha

import (
	"bytes"
	"crypto/rand"
	"github.com/dchest/uniuri"
	"io"
	"os"
)

// Standard number of digits in captcha.
const StdLength = 6

var globalStore = newStore()

// randomDigits return a byte slice of the given length containing random
// digits in range 0-9.
func randomDigits(length int) []byte {
	d := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, d); err != nil {
		panic(err)
	}
	for i := range d {
		d[i] %= 10
	}
	return d
}

// New creates a new captcha of the given length, saves it in the internal
// storage, and returns its id.
func New(length int) (id string) {
	id = uniuri.New()
	globalStore.saveCaptcha(id, randomDigits(length))
	return
}

// Reload generates and remembers new digits for the given captcha id.  This
// function returns false if there is no captcha with the given id.
//
// After calling this function, the image or audio presented to a user must be
// refreshed to show the new captcha representation (WriteImage and WriteAudio
// will write the new one).
func Reload(id string) bool {
	old := globalStore.getDigits(id)
	if old == nil {
		return false
	}
	globalStore.saveCaptcha(id, randomDigits(len(old)))
	return true
}

// WriteImage writes PNG-encoded image representation of the captcha with the
// given id. The image will have the given width and height.
func WriteImage(w io.Writer, id string, width, height int) os.Error {
	d := globalStore.getDigits(id)
	if d == nil {
		return os.NewError("captcha id not found")
	}
	_, err := NewImage(d, width, height).WriteTo(w)
	return err
}

// WriteAudio writes WAV-encoded audio representation of the captcha with the
// given id.
func WriteAudio(w io.Writer, id string) os.Error {
	d := globalStore.getDigits(id)
	if d == nil {
		return os.NewError("captcha id not found")
	}
	_, err := NewAudio(d).WriteTo(w)
	return err
}

// Verify returns true if the given digits are the ones that were used to
// create the given captcha id.
// 
// The function deletes the captcha with the given id from the internal
// storage, so that the same captcha can't be verified anymore.
func Verify(id string, digits []byte) bool {
	reald := globalStore.getDigitsClear(id)
	if reald == nil {
		return false
	}
	return bytes.Equal(digits, reald)
}

// Collect deletes expired and used captchas from the internal
// storage. It is called automatically by New function every CollectNum
// generated captchas, but still exported to enable freeing memory manually if
// needed.
//
// Collection is launched in a new goroutine, so this function returns
// immediately.
func Collect() {
	go globalStore.collect()
}
