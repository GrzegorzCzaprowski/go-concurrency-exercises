//////////////////////////////////////////////////////////////////////
//
// Given is some code to cache key-value pairs from a database into
// the main memory (to reduce access time). Note that golang's map are
// not entirely thread safe. Multiple readers are fine, but multiple
// writers are not. Change the code to make this thread safe.
//

package main

import "container/list"

// CacheSize determines how big the cache can grow
const CacheSize = 100

// KeyStoreCacheLoader is an interface for the KeyStoreCache
type KeyStoreCacheLoader interface {
	// Load implements a function where the cache should gets it's content from
	Load(string) string
}

// KeyStoreCache is a LRU cache for string key-value pairs
type KeyStoreCache struct {
	cache map[string]string
	pages list.List
	load  func(string) string
	mutex mutex
}

type mutex chan struct{}

func (m mutex) Lock() {
	m <- struct{}{}
}
func (m mutex) Unlock() {
	<-m
}

// New creates a new KeyStoreCache
func New(load KeyStoreCacheLoader) *KeyStoreCache {
	return &KeyStoreCache{
		load:  load.Load,
		cache: make(map[string]string),
	}
}

// Get gets the key from cache, loads it from the source if needed
func (k *KeyStoreCache) Get(key string) string {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	val, ok := k.cache[key]

	// Miss - load from database and save it in cache
	if !ok {
		val = k.load(key)
		k.cache[key] = val
		k.pages.PushFront(key)

		// if cache is full remove the least used item
		if len(k.cache) > CacheSize {
			delete(k.cache, k.pages.Back().Value.(string))
			k.pages.Remove(k.pages.Back())
		}
	}

	return val
}

// Loader implements KeyStoreLoader
type Loader struct {
	DB *MockDB
}

// Load gets the data from the database
func (l *Loader) Load(key string) string {
	val, err := l.DB.Get(key)
	if err != nil {
		panic(err)
	}

	return val
}

func run() *KeyStoreCache {
	loader := Loader{
		DB: GetMockDB(),
	}
	cache := New(&loader)

	RunMockServer(cache)

	return cache
}

func main() {
	run()
}
