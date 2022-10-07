package caches

import (
	"sync/atomic"
	"time"

	"cache-server/helpers"
)

const (
	// NeverDie is a const. if ttl=0, never die
	NeverDie = 0
)


// value is a struct packs data
type value struct {
	// data storages real data
	Data []byte

	//  ttl represents the life of this data
	// unit: second
	Ttl int64

	// ctime represents the time when this data is created
	Ctime int64
}


// newValue returns a data after packing
func newValue(data []byte, ttl int64) *value {
	return &value{
		// use copy to keep this data separated
		Data:	helpers.Copy(data),
		Ttl:	ttl,
		Ctime: 	time.Now().Unix(),
	}
}


// alive returns whether this data is alive
func (v *value) alive() bool {
	// First, determine whether there is an expiration date, 
	// and then determine whether the current time exceeds the expiration date of the data
	return v.Ttl == NeverDie || time.Now().Unix()-v.Ctime < v.Ttl
}

// visit returns the data
func (v *value) visit() []byte {
	atomic.SwapInt64(&v.Ctime, time.Now().Unix())
	return v.Data
}


















