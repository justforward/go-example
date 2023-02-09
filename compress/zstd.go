package compress

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

type Zstd struct {
}

func (gzip *Zstd) Compress(r io.Reader) (io.Reader, error) {

	pr, pw := io.Pipe()

	enc, err := zstd.NewWriter(pw, zstd.WithEncoderLevel(zstd.SpeedFastest))

	if err != nil {
		return nil, err
	}

	go func() {
		_, _ = io.Copy(enc, r)
		enc.Flush()
		enc.Close()
		pw.Close()
	}()

	return &StatReader{Reader: pr, Name: "zstd compress"}, nil

}

func (gzip *Zstd) Uncompress(r io.Reader) (io.Reader, error) {
	//pr, pw := io.Pipe()

	reader, err := zstd.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &StatReader{Reader: reader, Name: "zstd uncompress"}, nil

	//go func() {
	//	io.Copy(pw, r)
	//	dec.Close()
	//	pw.Close()
	//}()
	//
	//return &StatReader{Reader: pr, Name: "zstd uncompress"}, nil
}
