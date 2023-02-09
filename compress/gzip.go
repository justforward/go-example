package compress

import (
	zip "compress/gzip"
	"io"
)

type Gzip struct {
}

func (gzip *Gzip) Compress(r io.Reader) (io.Reader, error) {
	pr, pw := io.Pipe()

	enc, err := zip.NewWriterLevel(pw, zip.BestCompression)
	if err != nil {
		return nil, err
	}

	go func() {
		_, _ = io.Copy(enc, r)
		enc.Flush()
		enc.Close()
		pw.Close()
	}()

	return &StatReader{Reader: pr, Name: "gzip compress"}, nil
}

// Uncompress from reader
func (gzip *Gzip) Uncompress(r io.Reader) (io.Reader, error) {
	//pr, pw := io.Pipe()

	reader, err := zip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &StatReader{Reader: reader, Name: "gzip uncompress"}, nil

	//go func() {
	//	io.Copy(pw, r)
	//	dec.Close()
	//	pw.Close()
	//}()
	//
	//return &StatReader{Reader: pr, Name: "gzip uncompress"}, nil
}
