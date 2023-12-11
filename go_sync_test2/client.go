package main

import (
	"crypto/md5"
	"fmt"
	gosync "github.com/Redundancy/go-sync"
	"github.com/Redundancy/go-sync/blocksources"
	"github.com/Redundancy/go-sync/filechecksum"
	"github.com/Redundancy/go-sync/indexbuilder"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	BlockSize   = uint(1048576)
	Concurrency = 10
)

type HttpFile struct {
	curr int64
	name string

	hc       *http.Client
	statOnce sync.Once
	statSize int64
}

func (f *HttpFile) Read(p []byte) (n int, err error) {
	defer func() { f.curr += int64(n) }()

	// 发送一个http的请求
	req, err := http.NewRequest(http.MethodGet, f.buildURL(), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept-Ranges", "bytes")
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", f.curr, f.curr+int64(len(p))-1))
	resp, err := f.hc.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.Body.Read(p)
}

func (f *HttpFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.curr = offset
	case io.SeekCurrent:
		f.curr += offset
	case io.SeekEnd:
		size, err := f.stat()
		if err != nil {
			return 0, err
		}
		f.curr = size + offset
	}
	return f.curr, nil
}

func (f *HttpFile) ReadAt(p []byte, off int64) (n int, err error) {
	req, err := http.NewRequest(http.MethodGet, f.buildURL(), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept-Ranges", "bytes")
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", off, off+int64(len(p))-1))
	resp, err := f.hc.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.Body.Read(p)
}

func (f *HttpFile) stat() (int64, error) {
	var statErr error
	f.statOnce.Do(func() {
		var req *http.Request
		req, err := http.NewRequest(http.MethodOptions, f.buildURL(), nil)
		if err != nil {
			statErr = err
			return
		}

		resp, err := f.hc.Do(req)
		if err != nil {
			statErr = err
			return
		}
		defer func() { _ = resp.Body.Close() }()

		f.statSize, statErr = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	})

	return f.statSize, statErr
}

func (f *HttpFile) buildURL() string {
	return f.buildPartURL("/fs")
}

func (f *HttpFile) buildPartURL(part string) string {
	return "http://10.0.6.189:8080" + path.Join(part, f.name)
}

type HttpFileChecksumLookup struct {
	hf *HttpFile
}

func (l *HttpFileChecksumLookup) GetStrongChecksumForBlock(blockID int) []byte {
	req, err := http.NewRequest(http.MethodGet, l.hf.buildPartURL("/md5"), nil)
	if err != nil {
		panic(err)
	}

	//设定压缩？
	req.Header.Set("X-Checksum-Offset", strconv.Itoa(blockID*int(BlockSize)))
	req.Header.Set("X-Checksum-Length", strconv.Itoa(int(BlockSize)))
	resp, err := l.hf.hc.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	checksum, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return checksum
}

func (f *HttpFile) Summary() (gosync.FileSummary, error) {
	size, err := f.stat()
	if err != nil {
		panic(err)
	}

	g := filechecksum.NewFileChecksumGenerator(BlockSize)
	_, index, lookup, err := indexbuilder.BuildChecksumIndex(g, &HttpFile{name: f.name, hc: f.hc})
	if err != nil {
		return nil, err
	}

	return &gosync.BasicSummary{
		BlockSize:      BlockSize, // 分块得到数据 分块的数据使用udp的并发得到？
		BlockCount:     (uint(size) + BlockSize - 1) / BlockSize,
		FileSize:       size,
		ChecksumIndex:  index,
		ChecksumLookup: lookup,
		//ChecksumLookup: &HttpFileChecksumLookup{hf: f},
	}, nil
}

func main() {
	// gsync remote local
	if len(os.Args) != 3 {
		_, _ = fmt.Fprintf(os.Stderr, "%s remote local\n", os.Args[0])
		return
	}

	remote, local := os.Args[1], os.Args[2]
	if err := os.MkdirAll(filepath.Dir(local), 0755); err != nil {
		panic(err)
	}

	dstFile, err := os.OpenFile(local, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = dstFile.Close() }()

	hf := &HttpFile{name: remote, hc: http.DefaultClient}
	// summary
	summary, err := hf.Summary()
	if err != nil {
		panic(err)
	}

	resolver := blocksources.MakeFileSizedBlockResolver(
		uint64(summary.GetBlockSize()),
		summary.GetFileSize(),
	)

	rsync := gosync.RSync{
		Input:  hf,
		Output: dstFile,
		//OnClose: []io.Closer{dstFile},
		Source: blocksources.NewHttpBlockSource(hf.buildURL(), Concurrency, resolver, &filechecksum.HashVerifier{
			Hash:                md5.New(),
			BlockSize:           summary.GetBlockSize(),
			BlockChecksumGetter: summary,
		}),
		Summary: summary,
	}

	if err = rsync.Patch(); err != nil {
		panic(err)
	}
}
