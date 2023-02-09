package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"

	"github.com/klauspost/compress/zstd"
	//"github.com/pkg/profile"
)

const (
	ServerListenAddr = "127.0.0.1:8235"
	HttpServerAddr   = "http://127.0.0.1:8235/chunks"
)

func ReadMemStat(name string) {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	// Alloc is bytes of allocated heap objects.
	// TotalAlloc is cumulative bytes allocated for heap objects.
	// HeapAlloc is bytes of allocated heap objects.
	// Mallocs is the cumulative count of heap objects allocated.
	log.Printf("[%s]\nAlloc: %.5fMiB\nTotalAlloc: %.5fMiB\nHeadAlloc: %.5fMiB\nMallocs: %d\n\n", name,
		float64(memStat.Alloc)/1024.0/1024.0,
		float64(memStat.TotalAlloc)/1024.0/1024.0,
		float64(memStat.HeapAlloc)/1024.0/1024.0,
		memStat.Mallocs,
	)

}


func main() {
	//defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()
	//defer ReadMemStat("End")

	//ReadMemStat("Start")
	http.HandleFunc("/chunks", func(w http.ResponseWriter, r *http.Request) {
		rawReq, err := httputil.DumpRequest(r, false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Printf("[Server] Request:\n%s\n", rawReq)

		body, err := Decompress("Server", &StatReader{Reader: r.Body, Name: "[Server][Request Body]"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		compressBody, err := Compress("Server", body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		wn, we := io.Copy(w, compressBody)
		if we != nil {
			log.Printf("[Server] write to response error %q", we)
			return
		}
		log.Printf("[Server] write bytes = %d", wn)
	})

	go func() {
		if err := http.ListenAndServe(ServerListenAddr, nil); err != nil {
			panic(err)
		}
		log.Println("Http server shutdown ...")
	}()

	//for _, arg := range os.Args[1:] {
	filename := "/Users/mac/Downloads/Q3.sim"
	if err := SendCompressFile(filename); err != nil {
			panic(err)
		}
	//}

	var end runtime.MemStats
	runtime.ReadMemStats(&end)
}

func SendCompressFile(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = fp.Close() }()

	cr, err := Compress("Client", &StatReader{Reader: fp, Name: "[Client][File]"})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, HttpServerAddr, cr)
	if err != nil {
		return err
	}

	rawReq, err := httputil.DumpRequest(req, false)
	if err != nil {
		return err
	}
	log.Printf("[Client] Request:\n%s\n\n", rawReq)

	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	rawResp, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return err
	}
	log.Printf("[Client] Response:\n%s\n\n", rawResp)

	body, err := Decompress("Client", resp.Body)
	if err != nil {
		return err
	}

	bodySize, err := ReaderSize(body)
	if err != nil {
		return err
	}
	log.Printf("[Client] Response raw body size: %d", bodySize)

	return nil
}

type StatReader struct {
	io.Reader

	Name  string
	Reads int
}

func (r *StatReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	//log.Printf("[%s] read %d bytes with error %v", r.Name, n, err)

	r.Reads += n
	if err == io.EOF {
		log.Printf("%s total bytes = %d", r.Name, r.Reads)
	}

	return n, err
}

func Compress(name string, r io.Reader) (io.Reader, error) {
	pr, pw := io.Pipe()

	enc, err := zstd.NewWriter(pw)
	if err != nil {
		return nil, err
	}

	go func() {
		_, _ = io.Copy(enc, r)
		_ = enc.Flush()
		_ = enc.Close()
		_ = pw.Close()
	}()

	return &StatReader{Reader: pr, Name: fmt.Sprintf("[%s][ZStd_Compress]", name)}, nil
}

func Decompress(name string, r io.Reader) (io.Reader, error) {
	pr, pw := io.Pipe()

	dec, err := zstd.NewReader(pr)
	if err != nil {
		return nil, err
	}

	go func() {
		_, _ = io.Copy(pw, r)
		dec.Close()
		_ = pw.Close()
	}()

	return &StatReader{Reader: pr, Name: fmt.Sprintf("[%s][ZStd_Decompress]", name)}, nil
}

func ReaderSize(r io.Reader) (int64, error) {
	return io.Copy(io.Discard, r)
}
