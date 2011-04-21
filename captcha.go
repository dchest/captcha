package captcha

import (
	"bytes"
	"os"
	"rand"
	"time"
	crand "crypto/rand"
	"github.com/dchest/uniuri"
	"io"
	"container/list"
	"sync"
)

const (
	// Expiration time for captchas
	Expiration = 2 * 60 // 2 minutes
	// The number of captchas created that triggers garbage collection
	CollectNum = 100
	// The number of numbers to use in captcha
	NumCount = 6
)

// expValue stores timestamp and id of captchas. It is used in a list inside
// storage for indexing generated captchas by timestamp to enable garbage
// collection of expired captchas.
type expValue struct {
	timestamp int64
	id        string
}

// storage is an internal storage for captcha ids and their values.
type storage struct {
	mu  sync.RWMutex
	ids map[string][]byte
	exp *list.List
	// Number of items stored after last collection
	colNum int
}

func newStore() *storage {
	s := new(storage)
	s.ids = make(map[string][]byte)
	s.exp = list.New()
	return s
}

var store = newStore()

func init() {
	rand.Seed(time.Seconds())
}

func randomNumbers() []byte {
	n := make([]byte, NumCount)
	if _, err := io.ReadFull(crand.Reader, n); err != nil {
		panic(err)
	}
	for i := range n {
		n[i] %= 10
	}
	return n
}

// New creates a new captcha, saves it in the internal storage, and returns its
// id.
func New() string {
	ns := randomNumbers()
	id := uniuri.New()
	store.mu.Lock()
	defer store.mu.Unlock()
	store.ids[id] = ns
	store.exp.PushBack(expValue{time.Seconds(), id})
	store.colNum++
	if store.colNum > CollectNum {
		Collect()
		store.colNum = 0
	}
	return id
}

// WriteImage writes PNG-encoded captcha image of the given width and height
// with the given captcha id into the io.Writer.
func WriteImage(w io.Writer, id string, width, height int) os.Error {
	store.mu.RLock()
	defer store.mu.RUnlock()
	ns, ok := store.ids[id]
	if !ok {
		return os.NewError("captcha id not found")
	}
	return NewImage(ns, width, height).PNGEncode(w)
}

// Verify returns true if the given numbers are the numbers that were used to
// create the given captcha id.
// 
// The function deletes the captcha with the given id from the internal
// storage, so that the same captcha can't be used anymore.
func Verify(id string, numbers []byte) bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	realns, ok := store.ids[id]
	if !ok {
		return false
	}
	store.ids[id] = nil, false
	return bytes.Equal(numbers, realns)
}

// Collect deletes expired and used captchas from the internal
// storage. It is called automatically by New function every CollectNum
// generated captchas, but still exported to enable freeing memory manually if
// needed.
func Collect() {
	now := time.Seconds()
	store.mu.Lock()
	defer store.mu.Unlock()
	for e := store.exp.Front(); e != nil; e = e.Next() {
		ev, ok := e.Value.(expValue)
		if !ok {
			return
		}
		if ev.timestamp+Expiration < now {
			store.ids[ev.id] = nil, false
			store.exp.Remove(e)
		} else {
			return
		}
	}
}
