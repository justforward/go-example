package compress

import "io"

type Compressor interface {
	// Compress to writer, call Close() to flush content
	Compress(reader io.Reader) (io.Reader, error)
	// Uncompress from reader
	Uncompress(io.Reader) (io.Reader, error)
}
