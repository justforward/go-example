package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	// Set up the HTTP server
	http.HandleFunc("/file/download", serveFile)
	//http.HandleFunc("/file/stat", stat)
	http.ListenAndServe(":8081", nil)
}

func stat(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	filePath := params.Get("path")
	fd, err := os.Open(filePath)
	if err != nil {
		return
	}
	info, err := fd.Stat()
	if err != nil {
		return
	}

	// 构建一个结构体
	fileInfo := &File{
		Name:  info.Name(),
		IsDir: info.IsDir(),
		Size:  info.Size(),
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, &fileInfo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	w.Write(buf.Bytes())

}

type File struct {
	Name  string
	IsDir bool
	Size  int64
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	// Get the file path from the query string
	// 得到filepath
	//filePath := r.URL.Query().Get("path")

	filePath := "/Users/mac/test/test_append/Q3.sim"
	// Open the file to serve
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get the file info
	info, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Length header
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	fileExt := filepath.Ext(filePath)
	w.Header().Set("Content-Type", mime.TypeByExtension(fileExt))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	//w.WriteHeader(http.StatusPartialContent)
	http.ServeContent(w, r, filepath.Base(filePath), info.ModTime(), file)
}
