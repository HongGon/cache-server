package caches

import (
	"errors"
	"fmt"
	"sync"
	"time"
	// "cache-server/helpers"
)

// Cache is a struct to pack buffer
type Cache struct {
	// data is a map to storage all data
	// type of value is []byte for network transmission
	data map[string]*value

	options Options

	status *Status

	// lock to keep security
	lock *sync.RWMutex
}

// NewCache  return a Cache obj
func NewCache() *Cache {
	return NewCacheWith(DefaultOptions())
}

func NewCacheWith(options Options) *Cache {
	if cache, ok := recoverFromDumpFile(options.DumpFile); ok {
		return cache
	}
	return &Cache{
		data:    make(map[string]*value, 256),
        options: options,
        status:  newStatus(),
        lock:    &sync.RWMutex{},
	}
}

// recoverFromDumpFile recover cache from dumpfile
func recoverFromDumpFile(dumpFile string) (*Cache, bool) {
	cache, err := newEmptyDump().from(dumpFile)
	if err != nil {
		return nil, false
	}
	return cache, true
}


// Get return value by key, if not found return false
func (c *Cache) Get(key string) ([]byte, bool) {
	// use read-lock
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, false
	}

	//  if data is not alive, del it 
	// 注意这边对锁的操作，由于一开始加的是读锁，无法保证写的并发安全，而删除需要加写锁，读锁和写锁又是互斥的
	// 所以先将读锁释放，再上写锁，删除操作里面会加写锁，删除完之后，写锁释放，我们再上读锁
	if !value.alive() {
		c.lock.RUnlock()
		c.Delete(key)
		c.lock.RLock()
		return nil, false
	}
	return value.visit(), true
}

// Set save key-value to Cache
func (c *Cache) Set(key string, value []byte) error {
	return c.SetWithTTL(key, value , NeverDie)
}

// SetWithTTL add a k-v to cache
func (c *Cache) SetWithTTL(key string, value []byte, ttl int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if oldValue, ok := c.data[key]; ok {
		// if key exists, del it firstly
		c.status.subEntry(key, oldValue.Data)
	}
	
	if !c.checkEntrySize(key, value) {
		if oldValue, ok := c.data[key]; ok {
			c.status.addEntry(key, oldValue.Data)
		}
		return errors.New("the entry size will exceed if you set this entry")
	}
	// add a new k-v
	c.status.addEntry(key, value)
	c.data[key] = newValue(value, ttl)
	return nil
}




// Delete del key-value
func (c *Cache) Delete(key string) {
	// Delete can change status, use write lock
	c.lock.Lock()
	defer c.lock.Unlock()

	if oldValue, ok := c.data[key]; ok {
		// if this k-v exists
		c.status.subEntry(key, oldValue.Data)
		delete(c.data, key)
	}
}

// status return info of cache
func (c *Cache) Status() Status {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return *c.status
}

// checkEntrySize
func (c *Cache) checkEntrySize(newKey string, newValue []byte) bool {
	return c.status.entrySize()+int64(len(newKey))+int64(len(newValue)) <= c.options.MaxEntrySize*1024*1024
}


// gc
func (c *Cache) gc() {
	c.lock.Lock()
	defer c.lock.Unlock()
	// use count to record the count
	count := 0
	for key, value := range c.data {
		if !value.alive() {
			c.status.subEntry(key, value.Data)
			delete(c.data, key)
			count++
			if count >= c.options.MaxGcCount {
				break
			}
		}
	}
}

// AutoGC will start a GC timer
func (c *Cache) AutoGc() {
	go func ()  {
		ticker := time.NewTicker(time.Duration(c.options.GcDuration) * time.Minute)
		for {
			// use select to judge if reach the tick
			select {
			case <- ticker.C:
				c.gc()
			}
		}
	}()
}

// dump persist cache
func (c *Cache) dump() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	
	return newDump(c).to(c.options.DumpFile)
}

func (c *Cache) AutoDump() {
	go func() {
		ticker := time.NewTicker(time.Duration(c.options.DumpDuration)*time.Minute)
		for {
			select {
			case <- ticker.C:
				c.dump()
			}
		}
	}()
}
























