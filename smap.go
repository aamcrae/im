package im

import (
	"hash"
)

type sentry struct {
	hash uint64
	key string
	value any
}

type smap struct {
	hash hash.Hash64
	mod uint64
	tab [][]sentry
}

func (s smap) Get(key string) (any, bool) {
	s.hash.Reset()
	s.hash.Write([]byte(key))
	h := s.hash.Sum64()
	for _, ent := range s.tab[h % s.mod] {
		if ent.hash == h && ent.key == key {
			return ent.value, true
		}
	}
	return nil, false
}
