package caches

import (
	"sync"
	"cache-server/helpers"
)

// Cache is a struct to pack buffer
type Cache struct {
	// data is a map to storage all data
	// type of value is []byte for network transmission
	data map[string][]byte
	
	// count : the number of key-value
	count int64

	// lock to keep security
	lock *sync.RWMutex
}

// NewCache  return a Cache obj
func NewCache() *Cache {
	return &Cache{
		//256 is specified when creating the map
		//This 256 does not mean that only 256 key value pairs can be stored, but 256 slots are allocated in advance
		//Try to avoid map capacity expansion due to insufficient capacity, which will allocate memory and affect performance
		//Another point is that the fewer slots, the greater the probability of hash conflict, and the performance of map lookup will decline
		//256 here is not necessarily the best value, but depends on the actual situation
		data: make(map[string][]byte, 256),
		count: 0,
		lock: &sync.RWMutex{},
	}
}

// Get return value by key, if not found return false
func (c *Cache) Get(key string) ([]byte, bool) {
	// use read-lock
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

// Set save key-value to Cache
func (c *Cache) Set(key string, value []byte) {
	// Set operator can change the status of data, use write lock
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.data[key]
	if !ok {
		c.count++
	}
	// this method can copy value
	c.data[key] = helpers.Copy(value)
}


// Delete del key-value
func (c *Cache) Delete(key string) {
	// Delete can change status, use write lock
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.data[key]
	if ok {
		c.count--
		delete(c.data, key)
	}
}

// Count return the num of k-v
func (c *Cache) Count() int64 {
	// use read-lock
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.count
}



















