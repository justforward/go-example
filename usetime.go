package main

import (
	"io"
)

type TimeSpeed struct {
	io.Reader
	io.Writer
}

func (f TimeSpeed) Read(b []byte) (n int, err error) {
	return 0, nil
}

func r() {

}

func main() {

}
