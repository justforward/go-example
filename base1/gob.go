package main

import (
	"bytes"
	"encoding/gob"
	"os"
)

func main() {

	file, err := os.OpenFile("/Users/tan/test.hit", os.O_CREATE, 0646)
	var buffer bytes.Buffer
	enconder := gob.NewEncoder(&buffer)
	enconder.Encode()
}
