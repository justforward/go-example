package main

import (
	"bytes"
	"errors"
	"os"
	"strconv"

	"fmt"
	"go-example/compress"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	ioutil2 "github.com/coreos/etcd/pkg/ioutil"
)

const (
	ServerListenAddrHost = "127.0.0.1:8235"
	HttpServerAddrHost   = "http://127.0.0.1:8235"
)

func main() {

	http.HandleFunc("/chunks", func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		filePath := query.Get("path")
		offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
		if err != nil {

		}

		length, err := strconv.ParseInt(query.Get("length"), 10, 64)
		if err != nil {

		}

		compress := query.Get("compress")

		compressor := chooseCompressor(compress)

		// 从本地文件读取数据

		//log.Printf("[Server] Request:\n%s\n", rawReq)

		reader, err := ReadAt(filePath, offset, length)
		if err != nil {

		}

		defer reader.Close()

		compressReader, err := compressor.Compress(reader)
		if err != nil {
			return
		}

		wn, we := io.Copy(w, compressReader)
		if we != nil {
			log.Printf("[Server] write to response error %q", we)
			return
		}
		log.Printf("[Server] write bytes = %d", wn)
	})

	go func() {
		if err := http.ListenAndServe(ServerListenAddrHost, nil); err != nil {
			panic(err)
		}
		log.Println("Http server shutdown ...")
	}()

	fileName := "/Users/mac/Downloads/Q3.sim"

	offset := int64(0)
	length := int64(104857600)

	index := 0

	for index < 10 {
		compress := "zstd"
		if err := client(fileName, offset, length, compress); err != nil {
			panic(err)
		}

		offset = length
		length += int64(104857600)

		index++
	}

}

// 客户端请求
func client(fileName string, offset, length int64, compress string) error {

	compressor := chooseCompressor(compress)

	query := url.Values{}
	query.Set("path", fileName)
	query.Set("offset", fmt.Sprintf("%d", offset))
	query.Set("length", fmt.Sprintf("%d", length))
	query.Set("compress", compress)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/chunks?%s", HttpServerAddrHost, query.Encode()),
		bytes.NewReader(nil))

	if err != nil {
		return err
	}

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

	body, err := compressor.Uncompress(resp.Body)
	if err != nil {
		return err
	}

	bodySize, err := ReaderSize1(body)
	if err != nil {
		return err
	}
	log.Printf("[Client] Response raw body size: %d", bodySize)

	return nil

}

func chooseCompressor(fn string) compress.Compressor {

	if fn == "gzip" {
		return &compress.Gzip{}
	}

	if fn == "zstd" {
		return &compress.Zstd{}
	}

	return nil
}

// TODO: pool fd
func ReadAt(filePath string, offset, length int64) (io.ReadCloser, error) {
	// if file is symlink, return a empty reader
	info, err := os.Lstat(filePath)
	if err != nil {
		return nil, err
	}

	// open it
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		f.Close()
		return nil, nil
	}

	// seek it
	if err := seek(f, offset); err != nil {
		return nil, err
	}

	reader := (io.Reader)(f)
	if length >= 0 {
		reader = io.LimitReader(f, length)
	}

	// read it
	return ioutil2.ReaderAndCloser{
		Reader: reader,
		Closer: f,
	}, nil
}

func seek(f *os.File, offset int64) error {
	r, err := f.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	if r != int64(offset) {
		return errors.New("seeked but offset not same")
	}
	return nil
}

func ReaderSize1(r io.Reader) (int64, error) {
	return io.Copy(io.Discard, r)
}
