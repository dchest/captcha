package captcha

import (
	"container/list"
	"sync"
	"time"
)

const (
	// Expiration time for captchas
	Expiration = 2 * 60 // 2 minutes
	// The number of captchas created that triggers garbage collection
	CollectNum = 100
)

// expValue stores timestamp and id of captchas. It is used in a list inside
// store for indexing generated captchas by timestamp to enable garbage
// collection of expired captchas.
type expValue struct {
	timestamp int64
	id        string
}

// store is an internal store for captcha ids and their values.
type store struct {
	mu  sync.RWMutex
	ids map[string][]byte
	exp *list.List
	// Number of items stored after last collection
	colNum int
}

// newStore initializes and returns a new store.
func newStore() *store {
	s := new(store)
	s.ids = make(map[string][]byte)
	s.exp = list.New()
	return s
}

// saveCaptcha saves the captcha id and the corresponding digits.
func (s *store) saveCaptcha(id string, digits []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids[id] = digits
	s.exp.PushBack(expValue{time.Seconds(), id})
	s.colNum++
	if s.colNum > CollectNum {
		go s.collect()
		s.colNum = 0
	}
}

// getDigits returns the digits for the given id.
func (s *store) getDigits(id string) (digits []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	digits, _ = s.ids[id]
	return
}

// getDigitsClear returns the digits for the given id, and removes them from
// the store.
func (s *store) getDigitsClear(id string) (digits []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	digits, ok := s.ids[id]
	if !ok {
		return
	}
	s.ids[id] = nil, false
	// XXX(dchest) Index (s.exp) will be cleaned when collecting expired
	// captchas.  Can't clean it here, because we don't store reference to
	// expValue in the map. Maybe store it?
	return
}

// collect deletes expired captchas from the store.
func (s *store) collect() {
	now := time.Seconds()
	s.mu.Lock()
	defer s.mu.Unlock()
	for e := s.exp.Front(); e != nil; e = e.Next() {
		ev, ok := e.Value.(expValue)
		if !ok {
			return
		}
		if ev.timestamp+Expiration < now {
			s.ids[ev.id] = nil, false
			s.exp.Remove(e)
		} else {
			return
		}
	}
}
