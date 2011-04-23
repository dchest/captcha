package captcha

import (
	"bytes"
	"crypto/rand"
	"github.com/dchest/uniuri"
	"io"
	"os"
)

// Standard number of numbers in captcha
const StdLength = 6

var globalStore = newStore()

// randomNumbers return a byte slice of the given length containing random
// numbers in range 0-9.
func randomNumbers(length int) []byte {
	n := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, n); err != nil {
		panic(err)
	}
	for i := range n {
		n[i] %= 10
	}
	return n
}

// New creates a new captcha of the given length, saves it in the internal
// storage, and returns its id.
func New(length int) (id string) {
	id = uniuri.New()
	globalStore.saveCaptcha(id, randomNumbers(length))
	return
}

// Reload generates and remembers new numbers for the given captcha id.  This
// function returns false if there is no captcha with the given id.
//
// After calling this function, the image or audio presented to a user must be
// refreshed to show the new captcha representation (WriteImage and WriteAudio
// will write the new one).
func Reload(id string) bool {
	oldns := globalStore.getNumbers(id)
	if oldns == nil {
		return false
	}
	globalStore.saveCaptcha(id, randomNumbers(len(oldns)))
	return true
}

// WriteImage writes PNG-encoded image representation of the captcha with the
// given id. The image will have the given width and height.
func WriteImage(w io.Writer, id string, width, height int) os.Error {
	ns := globalStore.getNumbers(id)
	if ns == nil {
		return os.NewError("captcha id not found")
	}
	_, err := NewImage(ns, width, height).WriteTo(w)
	return err
}

// WriteAudio writes WAV-encoded audio representation of the captcha with the
// given id.
func WriteAudio(w io.Writer, id string) os.Error {
	ns := globalStore.getNumbers(id)
	if ns == nil {
		return os.NewError("captcha id not found")
	}
	_, err := NewAudio(ns).WriteTo(w)
	return err
}

// Verify returns true if the given numbers are the numbers that were used to
// create the given captcha id.
// 
// The function deletes the captcha with the given id from the internal
// storage, so that the same captcha can't be verified anymore.
func Verify(id string, numbers []byte) bool {
	realns := globalStore.getNumbersClear(id)
	if realns == nil {
		return false
	}
	return bytes.Equal(numbers, realns)
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
