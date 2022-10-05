package caches


// options is struct of option

type Options struct {
	// MaxEntrySize is the max of write full
	MaxEntrySize	int64
	// MaxGcCount is the max of GC
	MaxGcCount		int
	// GcDuration is the duration of GC
	GcDuration		int64
}

// DefaultOptions returns a default option obj
func DefaultOptions() Options {
	return Options{
		MaxEntrySize: 	int64(4),	// default 4 GB
		MaxGcCount: 	1000,	// default 1000
		GcDuration: 	60,		// default an hour
	}
}








