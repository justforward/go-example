package fs

import (
	"io"
	"os"
)

type FileSystem interface {
	Ls(path string, pageOffset, pageSize int64) ([]os.FileInfo, error)
	Stat(path string) (os.FileInfo, error)
	Download(path string, beginOffset int64, endOffset int64) (io.ReadCloser, error)
}
