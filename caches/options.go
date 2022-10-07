package caches


// options is struct of option

type Options struct {

	// DumpFile is the persist file
	DumpFile string

	// Duration
	DumpDuration int64

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
		DumpFile: "cache-server.dump",
		DumpDuration: 30,
		MaxEntrySize: 	int64(4),	// default 4 GB
		MaxGcCount: 	1000,	// default 1000
		GcDuration: 	60,		// default an hour
	}
}








