package compress

import (
	"io"
	"log"
)

type StatReader struct {
	io.Reader

	Name  string
	Reads int
}

func (r *StatReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	log.Printf("[%s] read %d bytes with error %v", r.Name, n, err)

	r.Reads += n
	if err == io.EOF {
		log.Printf("%s total bytes = %d", r.Name, r.Reads)
	}

	return n, err
}
