package memoize

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
)

type memoizer1To1[I1 any, O1 any] struct {
	storage *cache.Cache
	group   singleflight.Group
	fn      func(I1) (O1, error)
}

// Memoize1To1 memoizes a function with 1 input and 1 output parameter. If the underlying function errors,
// the result is not cached.
func Memoize1To1[I1 any, O1 any](defaultExpiration time.Duration, fn func(I1) (O1, error)) func(I1) (O1, error) {
	m := &memoizer1To1[I1, O1]{
		storage: cache.New(defaultExpiration, defaultExpiration/2),
		group:   singleflight.Group{},
		fn:      fn,
	}
	return m.do
}

func (m *memoizer1To1[I1, O1]) do(i1 I1) (O1, error) {
	key := fmt.Sprint(i1)
	_r, err, _ := m.group.Do(key, func() (interface{}, error) {
		r, found := m.storage.Get(key)
		if found {
			return r, nil
		}

		r, err := m.fn(i1)
		if err != nil {
			// don't cache
			return r, err
		}

		m.storage.Set(key, r, cache.DefaultExpiration)
		return r, err
	})

	var r O1
	if _r != nil {
		r = _r.(O1)
	}
	return r, err
}
