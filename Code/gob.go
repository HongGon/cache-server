package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
)

func main() {
	m := map[string]int {
		"A":1,
		"B":2,
		"C":3,
	}
	file, err := ioutil.TempFile("","gob_test_*")
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(m)
	if err != nil {
		panic(err)
	}
	file.Seek(0,io.SeekStart)
	decoder := gob.NewDecoder(file)

	newM := map[string]int{}
	err = decoder.Decode(&newM)
	if err != nil {
		panic(err)
	}
	fmt.Println(newM)

}




