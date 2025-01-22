package util

import (
	"iter"
	"sync"
)

type kvpair[K comparable, T any] struct {
	key K
	val T
}

// RWMap is threads safe map.
type RWMap[K comparable, T any] struct {
	m   map[K]T
	mux sync.RWMutex
}

// Makes initial capacity, it can auto expand later.
func (rwm *RWMap[K, T]) Init(capacity int) {
	rwm.mux.Lock()
	defer rwm.mux.Unlock()
	rwm.m = make(map[K]T, capacity)
}

// Len returns number of key-value pairs.
func (rwm *RWMap[K, T]) Len() int {
	rwm.mux.RLock()
	defer rwm.mux.RUnlock()
	return len(rwm.m)
}

// Has detects whether pair is present.
func (rwm *RWMap[K, T]) Has(key K) (ok bool) {
	rwm.mux.RLock()
	defer rwm.mux.RUnlock()
	_, ok = rwm.m[key]
	return
}

// Get returns value by pointed key.
func (rwm *RWMap[K, T]) Get(key K) (ret T, ok bool) {
	rwm.mux.RLock()
	defer rwm.mux.RUnlock()
	ret, ok = rwm.m[key]
	return
}

// Set inserts given key-value pair.
func (rwm *RWMap[K, T]) Set(key K, val T) {
	rwm.mux.Lock()
	defer rwm.mux.Unlock()
	rwm.m[key] = val
}

// Delete removes pair from map.
func (rwm *RWMap[K, T]) Delete(key K) {
	rwm.mux.Lock()
	defer rwm.mux.Unlock()
	delete(rwm.m, key)
}

// GetAndDelete removes pair from map and returns the value if it was.
func (rwm *RWMap[K, T]) GetAndDelete(key K) (ret T, ok bool) {
	rwm.mux.Lock()
	defer rwm.mux.Unlock()
	if ret, ok = rwm.m[key]; ok {
		delete(rwm.m, key)
	}
	return
}

// Items returns iterator for all key-value pairs. Iterator makes copy of
// the state and then yields for each pair.
func (rwm *RWMap[K, T]) Items() iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		var buf []kvpair[K, T]
		func() {
			rwm.mux.RLock()
			defer rwm.mux.RUnlock()
			buf = make([]kvpair[K, T], len(rwm.m))
			var i int
			for k, v := range rwm.m {
				buf[i].key, buf[i].val = k, v
				i++
			}
		}() // unlock when copy is ready
		for _, pair := range buf {
			if !yield(pair.key, pair.val) {
				return
			}
		}
	}
}

// Cache is LRU & FIFO threads safe cache.
type Cache[K comparable, T any] struct {
	seq []kvpair[K, T] // sequence of key-value pairs
	idx map[K]int      // map with pairs positions pointed by keys
	efn func(K, T)     // exit function, called on pair remove
	mux sync.Mutex
}

// NewCache returns pointer to new Cache object.
func NewCache[K comparable, T any]() *Cache[K, T] {
	return &Cache[K, T]{
		idx: map[K]int{},
	}
}

// OnRemove changes callback function that is called when removing a pair.
func (c *Cache[K, T]) OnRemove(efn func(K, T)) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.efn = efn
}

// Len returns number of key-value pairs.
func (c *Cache[K, T]) Len() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return len(c.seq)
}

// Has detects whether pair is present.
func (c *Cache[K, T]) Has(key K) (ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok = c.idx[key]
	return
}

// Peek returns value pointed by given key.
func (c *Cache[K, T]) Peek(key K) (ret T, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	var n int
	if n, ok = c.idx[key]; ok {
		ret = c.seq[n].val
	}
	return
}

// Get returns value pointed by given key, and brings the pair to top of cache.
func (c *Cache[K, T]) Get(key K) (ret T, ok bool) {
	var n int

	c.mux.Lock()
	defer c.mux.Unlock()

	n, ok = c.idx[key]
	if ok {
		var pair = c.seq[n]
		ret = pair.val
		copy(c.seq[n:], c.seq[n+1:])
		c.seq[len(c.seq)-1] = pair
		for i := n; i < len(c.seq); i++ {
			c.idx[c.seq[i].key] = i
		}
	}
	return
}

func (c *Cache[K, T]) Poke(key K, val T) {
	c.mux.Lock()
	defer c.mux.Unlock()

	var n, ok = c.idx[key]
	if ok {
		c.seq[n].val = val
	} else {
		c.idx[key] = len(c.seq)
		c.seq = append(c.seq, kvpair[K, T]{
			key: key,
			val: val,
		})
	}
}

func (c *Cache[K, T]) Set(key K, val T) {
	c.mux.Lock()
	defer c.mux.Unlock()

	var n, ok = c.idx[key]
	if ok {
		var pair = c.seq[n]
		pair.val = val
		copy(c.seq[n:], c.seq[n+1:])
		c.seq[len(c.seq)-1] = pair
		for i := n; i < len(c.seq); i++ {
			c.idx[c.seq[i].key] = i
		}
	} else {
		c.idx[key] = len(c.seq)
		c.seq = append(c.seq, kvpair[K, T]{
			key: key,
			val: val,
		})
	}
}

func (c *Cache[K, T]) Remove(key K) (ok bool) {
	var n int

	c.mux.Lock()
	defer c.mux.Unlock()

	n, ok = c.idx[key]
	if ok {
		var pair = c.seq[n]
		if c.efn != nil {
			c.efn(pair.key, pair.val)
		}
		delete(c.idx, key)
		copy(c.seq[n:], c.seq[n+1:])
		c.seq = c.seq[:len(c.seq)-1]
		for i := n; i < len(c.seq); i++ {
			c.idx[c.seq[i].key] = i
		}
	}
	return
}

// Items returns iterator for all key-value pairs. Iterator makes copy of
// the state and then yields for each pair.
func (c *Cache[K, T]) Items() iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		c.mux.Lock()
		var s = append([]kvpair[K, T]{}, c.seq...) // make non-nil copy
		c.mux.Unlock()

		for _, pair := range s {
			if !yield(pair.key, pair.val) {
				return
			}
		}
	}
}

// Until removes first some entries from cache until given func returns true.
func (c *Cache[K, T]) Until(f func(K, T) bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	var n = 0
	if c.efn != nil {
		for _, pair := range c.seq {
			if f(pair.key, pair.val) {
				c.efn(pair.key, pair.val)
				delete(c.idx, pair.key)
				n++
			} else {
				break
			}
		}
	} else {
		for _, pair := range c.seq {
			if f(pair.key, pair.val) {
				delete(c.idx, pair.key)
				n++
			} else {
				break
			}
		}
	}
	c.seq = c.seq[n:]
}

// Free removes n first entries from cache.
func (c *Cache[K, T]) Free(n int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if n <= 0 {
		return
	}
	if n >= len(c.seq) {
		if c.efn != nil {
			for _, pair := range c.seq {
				c.efn(pair.key, pair.val)
			}
		}
		c.idx = map[K]int{}
		c.seq = nil
		return
	}

	if c.efn != nil {
		for i := 0; i < n; i++ {
			c.efn(c.seq[i].key, c.seq[i].val)
			delete(c.idx, c.seq[i].key)
		}
	} else {
		for i := 0; i < n; i++ {
			delete(c.idx, c.seq[i].key)
		}
	}
	c.seq = c.seq[n:]
}

// ToLimit brings cache to limited count of entries.
func (c *Cache[K, T]) ToLimit(limit int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if limit >= len(c.seq) {
		return
	}
	if limit <= 0 {
		if c.efn != nil {
			for _, pair := range c.seq {
				c.efn(pair.key, pair.val)
			}
		}
		c.idx = map[K]int{}
		c.seq = nil
		return
	}

	var n = len(c.seq) - limit
	if c.efn != nil {
		for i := 0; i < n; i++ {
			c.efn(c.seq[i].key, c.seq[i].val)
			delete(c.idx, c.seq[i].key)
		}
	} else {
		for i := 0; i < n; i++ {
			delete(c.idx, c.seq[i].key)
		}
	}
	c.seq = c.seq[n:]
}

// Sizer is interface that determine structure size itself.
type Sizer interface {
	Size() int64
}

// CacheSize returns size of given cache.
func CacheSize[K comparable, T Sizer](cache *Cache[K, T]) (size int64) {
	for _, val := range cache.Items() {
		size += val.Size()
	}
	return
}

// Bimap is bidirectional threads safe map.
type Bimap[K comparable, T comparable] struct {
	dir map[K]T // direct order
	rev map[T]K // reverse order
	mux sync.RWMutex
}

// NewBimap returns pointer to new Bimap object.
func NewBimap[K comparable, T comparable]() *Bimap[K, T] {
	return &Bimap[K, T]{
		dir: map[K]T{},
		rev: map[T]K{},
	}
}

// Len returns number of key-value pairs.
func (m *Bimap[K, T]) Len() int {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return len(m.dir)
}

// GetDir returns element in direct order, i.e. value pointed by key.
func (m *Bimap[K, T]) GetDir(key K) (val T, ok bool) {
	m.mux.RLock()
	val, ok = m.dir[key]
	m.mux.RUnlock()
	return
}

// GetRev returns element in reverse order, i.e. key pointed by value.
func (m *Bimap[K, T]) GetRev(val T) (key K, ok bool) {
	m.mux.RLock()
	key, ok = m.rev[val]
	m.mux.RUnlock()
	return
}

// Set inserts key-value pair into map.
func (m *Bimap[K, T]) Set(key K, val T) {
	m.mux.Lock()
	m.dir[key] = val
	m.rev[val] = key
	m.mux.Unlock()
}

// DeleteDir deletes key-value pair pointed by key, and returns deleted value.
func (m *Bimap[K, T]) DeleteDir(key K) (val T, ok bool) {
	m.mux.Lock()
	if val, ok = m.dir[key]; ok {
		delete(m.dir, key)
		delete(m.rev, val)
	}
	m.mux.Unlock()
	return
}

// DeleteRev deletes key-value pair pointed by value, and returns deleted key.
func (m *Bimap[K, T]) DeleteRev(val T) (key K, ok bool) {
	m.mux.Lock()
	if key, ok = m.rev[val]; ok {
		delete(m.dir, key)
		delete(m.rev, val)
	}
	m.mux.Unlock()
	return
}

// The End.
