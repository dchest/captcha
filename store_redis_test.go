// Contributed 2020 by Hari
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"bytes"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestRedisSetGet(t *testing.T) {
	s, err := NewRedisStore(&redis.Options{Addr: "localhost:6379", DB: 0}, 1*time.Minute, DefaultMaxRedisKeys, DefaultRedisPrefixKey)
	if err != nil {
		t.Errorf(err.Error())
	}
	id := "redis-id-no-clear"
	d := RandomDigits(10)
	s.Set(id, d)
	d2 := s.Get(id, false)
	if d2 == nil || !bytes.Equal(d, d2) {
		t.Errorf("saved %v, getDigits returned got %v", d, d2)
	}
}

func TestRedisGetClear(t *testing.T) {
	s, err := NewRedisStore(&redis.Options{Addr: "localhost:6379", DB: 0}, Expiration, DefaultMaxRedisKeys, DefaultRedisPrefixKey)
	if err != nil {
		t.Errorf(err.Error())
	}
	id := "redis-id"
	d := RandomDigits(10)
	s.Set(id, d)
	d2 := s.Get(id, true)
	if d2 == nil || !bytes.Equal(d, d2) {
		t.Errorf("saved %v, getDigits returned got %v", d, d2)
	}
	d2 = s.Get(id, false)
	if d2 != nil {
		t.Errorf("getDigitClear didn't clear (%q=%v)", id, d2)
	}
}

func BenchmarkRedisMaxKeys(b *testing.B) {
	maxKeys := 101

	b.StopTimer()
	d := RandomDigits(10)
	s, err := NewRedisStore(&redis.Options{Addr: "localhost:6379", DB: 0}, 1*time.Minute, int64(maxKeys), DefaultRedisPrefixKey)
	if err != nil {
		b.Errorf(err.Error())
	}
	ids := make([]string, maxKeys)
	for i := range ids {
		ids[i] = randomId()
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < maxKeys; j++ {
			s.Set(ids[j], d)
		}
	}
}
