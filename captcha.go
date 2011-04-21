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
	expiration = 2 * 60 // 2 minutes
	collectNum = 100    // number of items that triggers collection
)

type expValue struct {
	timestamp int64
	id        string
}

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
	n := make([]byte, 6)
	if _, err := io.ReadFull(crand.Reader, n); err != nil {
		panic(err)
	}
	for i := range n {
		n[i] %= 10
	}
	return n
}

func New() string {
	ns := randomNumbers()
	id := uniuri.New()
	store.mu.Lock()
	defer store.mu.Unlock()
	store.ids[id] = ns
	store.exp.PushBack(expValue{time.Seconds(), id})
	store.colNum++
	if store.colNum > collectNum {
		Collect()
		store.colNum = 0
	}
	return id
}

func WriteImage(w io.Writer, id string, width, height int) os.Error {
	store.mu.RLock()
	defer store.mu.RUnlock()
	ns, ok := store.ids[id]
	if !ok {
		return os.NewError("captcha id not found")
	}
	return NewImage(ns, width, height).PNGEncode(w)
}

func Verify(w io.Writer, id string, ns []byte) bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	realns, ok := store.ids[id]
	if !ok {
		return false
	}
	store.ids[id] = nil, false
	return bytes.Equal(ns, realns)
}

func Collect() {
	now := time.Seconds()
	store.mu.Lock()
	defer store.mu.Unlock()
	for e := store.exp.Front(); e != nil; e = e.Next() {
		ev, ok := e.Value.(expValue)
		if !ok {
			return
		}
		if ev.timestamp+expiration < now {
			store.ids[ev.id] = nil, false
			store.exp.Remove(e)
		} else {
			return
		}
	}
}
