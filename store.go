package captcha

import (
	"container/list"
	"sync"
	"time"
)

// An object implementing Store interface can be registered with SetCustomStore
// function to handle storage and retrieval of captcha ids and solutions for
// them, replacing the default memory store.
type Store interface {
	// Set sets the digits for the captcha id.
	Set(id string, digits []byte)

	// Get returns stored digits for the captcha id. Clear indicates
	// whether the captcha must be deleted from the store.
	Get(id string, clear bool) (digits []byte)

	// Collect deletes expired captchas from the store.  For custom stores
	// this method is not called automatically, it is only wired to the
	// package's Collect function.  Custom stores must implement their own
	// procedure for calling Collect, for example, in Set method.
	Collect()
}

// expValue stores timestamp and id of captchas. It is used in the list inside
// memoryStore for indexing generated captchas by timestamp to enable garbage
// collection of expired captchas.
type expValue struct {
	timestamp int64
	id        string
}

// memoryStore is an internal store for captcha ids and their values.
type memoryStore struct {
	mu  sync.RWMutex
	ids map[string][]byte
	exp *list.List
	// Number of items stored after last collection.
	numStored int
	// Number of saved items that triggers collection.
	collectNum int
	// Expiration time of captchas.
	expiration int64
}

// NewMemoryStore returns a new standard memory store for captchas with the
// given collection threshold and expiration time in seconds. The returned
// store must be registered with SetCustomStore to replace the default one.
func NewMemoryStore(collectNum int, expiration int64) Store {
	s := new(memoryStore)
	s.ids = make(map[string][]byte)
	s.exp = list.New()
	s.collectNum = collectNum
	s.expiration = expiration
	return s
}

func (s *memoryStore) Set(id string, digits []byte) {
	s.mu.Lock()
	s.ids[id] = digits
	s.exp.PushBack(expValue{time.Seconds(), id})
	s.numStored++
	s.mu.Unlock()
	if s.numStored > s.collectNum {
		go s.Collect()
	}
}

func (s *memoryStore) Get(id string, clear bool) (digits []byte) {
	if !clear {
		// When we don't need to clear captcha, acquire read lock.
		s.mu.RLock()
		defer s.mu.RUnlock()
	} else {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	digits, ok := s.ids[id]
	if !ok {
		return
	}
	if clear {
		s.ids[id] = nil, false
		// XXX(dchest) Index (s.exp) will be cleaned when collecting expired
		// captchas.  Can't clean it here, because we don't store reference to
		// expValue in the map. Maybe store it?
	}
	return
}

func (s *memoryStore) Collect() {
	now := time.Seconds()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.numStored = 0
	for e := s.exp.Front(); e != nil; {
		ev, ok := e.Value.(expValue)
		if !ok {
			return
		}
		if ev.timestamp+s.expiration < now {
			s.ids[ev.id] = nil, false
			next := e.Next()
			s.exp.Remove(e)
			e = next
		} else {
			return
		}
	}
}
