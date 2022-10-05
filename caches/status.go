package caches


// Status is a struct represents cache info

type Status struct {
	// Cound records the num of data
	Count int `json:"count`
	
	// KeySize records the capacity that the key occupy
	KeySize int64 `json:"keySize"`

	// Value Size records the capacity that value occupy
	ValueSize int64 `json:"valueSize"`
}

// newStatus returns a obj pointer of cache info
func newStatus() *Status {
	return &Status{
		Count:		0,
		KeySize: 	0,
		ValueSize:  0,
	}
}


// addEntry can record the info of key and value
func (s *Status) addEntry(key string, value []byte) {
	// add a k-v, count+1
	s.Count++
	s.KeySize += int64(len(key))
	s.ValueSize += int64(len(value))
}

// subEntry can remove the info of k-v
func (s *Status) subEntry(key string, value []byte) {
	s.Count--
	s.KeySize -= int64(len(key))
	s.ValueSize -= int64(len(value))
}

// entrySize returns the capacity of k-v
func (s *Status) entrySize() int64 {
	return s.KeySize + s.ValueSize
}




















