package caches

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"
)

type dump struct {
	// Data storages k-v
	Data map[string]*value

	// Options records config
	Options Options
	
	// Status records cache
	Status *Status
}

// newEmptyDump create a empty obj of dump
func newEmptyDump() *dump {
	return &dump{}
}

// newDump create a dump obj
func newDump(c *Cache) *dump {
	return &dump {
		Data: c.data,
		Options: c.options,
		Status: c.status,
	}
}


// newSuffix return a file name
func nowSuffix() string {
    return "." + time.Now().Format("20060102150405")
}

func (d *dump) to(dumpFile string) error {
	newDumpFile := dumpFile + nowSuffix()
	file, err := os.OpenFile(newDumpFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = gob.NewEncoder(file).Encode(d)
	if err != nil {
		file.Close()
		os.Remove(newDumpFile)
		return err
	}

	os.Remove(dumpFile)
	file.Close()
	return os.Rename(newDumpFile,dumpFile)
}

func (d *dump) from(dumpFile string) (*Cache,error) {
	// read dumpFile 
	file, err := os.Open(dumpFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err = gob.NewDecoder(file).Decode(d); err != nil {
		return nil, err
	}
	
	return &Cache{
		data:		d.Data,
		options:	d.Options,
		status:		d.Status,
		lock:		&sync.RWMutex{},
	}, nil
}





