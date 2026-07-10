package deduper

import (
	"context"
	"hash/fnv"
	"regexp"
	"strings"
	"sync"
)

var _ Deduper = (*hashmap)(nil)

var rHexId = regexp.MustCompile(`0x[0-9a-fA-F]+:0x[0-9a-fA-F]+`)

type hashmap struct {
	mux  *sync.RWMutex
	seen map[uint64]struct{}
}

func normalizeKey(key string) string {
	if strings.Contains(key, "google.com/maps") || strings.Contains(key, "maps.google") {
		if m := rHexId.FindString(key); m != "" {
			return strings.ToLower(m)
		}
		if idx := strings.Index(key, "/maps/place/"); idx != -1 {
			sub := key[idx+len("/maps/place/"):]
			if nextSlash := strings.IndexAny(sub, "/?@"); nextSlash != -1 {
				return strings.ToLower(sub[:nextSlash])
			}
			return strings.ToLower(sub)
		}
	}
	return strings.ToLower(strings.TrimSpace(key))
}

func (d *hashmap) AddIfNotExists(_ context.Context, key string) bool {
	normKey := normalizeKey(key)

	d.mux.RLock()
	if _, ok := d.seen[d.hash(normKey)]; ok {
		d.mux.RUnlock()
		return false
	}

	d.mux.RUnlock()

	d.mux.Lock()
	defer d.mux.Unlock()

	if _, ok := d.seen[d.hash(normKey)]; ok {
		return false
	}

	d.seen[d.hash(normKey)] = struct{}{}

	return true
}

func (d *hashmap) hash(key string) uint64 {
	h := fnv.New64()
	h.Write([]byte(key))

	return h.Sum64()
}
