package captcha

import (
	"bytes"
	"github.com/dchest/uniuri"
	"testing"
)

func TestSetGet(t *testing.T) {
	s := NewMemoryStore(CollectNum, Expiration)
	id := "captcha id"
	d := RandomDigits(10)
	s.Set(id, d)
	d2 := s.Get(id, false)
	if d2 == nil || !bytes.Equal(d, d2) {
		t.Errorf("saved %v, getDigits returned got %v", d, d2)
	}
}

func TestGetClear(t *testing.T) {
	s := NewMemoryStore(CollectNum, Expiration)
	id := "captcha id"
	d := RandomDigits(10)
	s.Set(id, d)
	d2 := s.Get(id, true)
	if d2 == nil || !bytes.Equal(d, d2) {
		t.Errorf("saved %v, getDigitsClear returned got %v", d, d2)
	}
	d2 = s.Get(id, false)
	if d2 != nil {
		t.Errorf("getDigitClear didn't clear (%q=%v)", id, d2)
	}
}

func TestCollect(t *testing.T) {
	//TODO(dchest): can't test automatic collection when saving, because
	//it's currently launched in a different goroutine.
	s := NewMemoryStore(10, -1)
	// create 10 ids
	ids := make([]string, 10)
	d := RandomDigits(10)
	for i := range ids {
		ids[i] = uniuri.New()
		s.Set(ids[i], d)
	}
	s.(*memoryStore).collect()
	// Must be already collected
	nc := 0
	for i := range ids {
		d2 := s.Get(ids[i], false)
		if d2 != nil {
			t.Errorf("%d: not collected", i)
			nc++
		}
	}
	if nc > 0 {
		t.Errorf("= not collected %d out of %d captchas", nc, len(ids))
	}
}
