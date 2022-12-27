package main

import (
	"fmt"
	"github.com/chmduquesne/rollinghash/buzhash64"
	"github.com/chmduquesne/rollinghash/rabinkarp64"
	"io"
	"os"
)

import (
	"github.com/chmduquesne/rollinghash"
)

var WindowSize = 64
var Pattern = uint64(8)<<20 - 1 // 4 ~ 8 MiB


//go:inline
func Match(v uint64) bool {
	return (v^Pattern)&Pattern == Pattern
}

func Chunk(r io.Reader, h rollinghash.Hash64, name string) error {
	// 读取64 位的文件
	buf := make([]byte, WindowSize)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	_, _ = h.Write(buf)
	if n != WindowSize {
		return nil
	}

	buf = make([]byte, 32768)
	for offset, boundary := 64, 0; true; {
		n, err = r.Read(buf)
		if err != nil {
			return nil
		}

		for _, b := range buf[:n] {
			h.Roll(b)

			// 如果最后几位是0的时候 当做一个界限
			if Match(h.Sum64()) {
				fmt.Printf("[%10s] in offset %10d rolling hash = %064b with chunk size is %03.10fM\n",
					name, offset, h.Sum64(), float64(offset-boundary)/1024/1024)
				boundary = offset
			}
			// 每次添加
			offset++
		}
	}

	return nil
}

func WithResettableReader(r io.ReadSeeker, h rollinghash.Hash64, name string) error {
	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	//
	return Chunk(r, h, name)
}


func main() {
	pol, err := rabinkarp64.RandomPolynomial(0x123456787654321)
	if err != nil {
		//panic(err)
	}
	fp, err := os.Open("./aaa.mp4")
	if err != nil {
		//panic(err)
	}
	defer func() { _ = fp.Close() }()

	// 进行文件的读写
	err = WithResettableReader(fp, rabinkarp64.NewFromPol(pol), "Rabin karp")
	if err != nil {
		//panic(err)
	}

	err = WithResettableReader(fp, buzhash64.New(), "Buz")
	if err != nil {
		//panic(err)
	}
}
