package custom

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {

	http.HandleFunc("/file/download", download)
	http.HandleFunc("/file/stat", stat)
	http.HandleFunc("/file/ls", ls)

	port := os.Args[1]
	http.ListenAndServe(":"+port, nil)
}

// range的请求
func download(w http.ResponseWriter, r *http.Request) {
	FilePath := r.URL.Query().Get("file_path")
	file, err := os.Open(FilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	fileExt := filepath.Ext(FilePath)
	w.Header().Set("Content-Type", mime.TypeByExtension(fileExt))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(FilePath)))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	http.ServeContent(w, r, filepath.Base(FilePath), info.ModTime(), file)
}

func stat(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	filePath := params.Get("file_path")
	fmt.Println(fmt.Sprintf("rsync_server get file path:%s", filePath))
	fd, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	info, err := fd.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("stat file name :%+v", info.Name()))

	fileInfo := &FileInfos{
		FileName:    info.Name(),
		FileSize:    info.Size(),
		FileModTime: info.ModTime(),
		FileIsDir:   info.IsDir(),
	}

	fmt.Println(fmt.Sprintf("stat file info:%+v", fileInfo))

	data, err := json.Marshal(fileInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	w.Write(data)
}

func ls(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	filePath := params.Get("path")
	pageOffset, err := strconv.Atoi(params.Get("pageOffset"))
	if err != nil {
		panic(err)
	}

	pageSize, err := strconv.Atoi(params.Get("pageSize"))
	if err != nil {
		panic(err)
	}

	infos, err := os.ReadDir(filePath)

	if err != nil {
		panic(err)
	}

	if len(infos) == 0 {
		panic(errors.New("file infos len is 0"))
	}

	start := pageOffset
	if pageOffset > 0 {
		start = pageSize * (pageOffset - 1)
	}

	if start >= len(infos) {
		panic(errors.New("start > len(infos)"))
	} else {
		end := start + pageSize
		if end > len(infos) {
			end = len(infos)
		}
		infos = infos[start:end]
	}

	var result []*FileInfos
	for _, info := range infos {
		fileInfo, err := info.Info()
		if err != nil {
			panic(err)
		}
		result = append(result, ToFileInfo(fileInfo))
	}

	var ans []os.FileInfo
	for _, res := range result {
		ans = append(ans, res)
	}

	data, err := json.Marshal(ans)
	if err != nil {
		panic(err)
	}

	// Set the Content-Type header
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON data to the response
	w.Write(data)
}

type FileInfos struct {
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	FileModTime time.Time `json:"file_mod_time"`
	FileIsDir   bool      `json:"file_is_dir"`
}

func ToFileInfo(info os.FileInfo) *FileInfos {
	return &FileInfos{
		FileName:    info.Name(),
		FileSize:    info.Size(),
		FileModTime: info.ModTime(),
		FileIsDir:   info.IsDir(),
	}
}

func (f *FileInfos) Name() string {
	return f.FileName
}

func (f *FileInfos) Size() int64 {
	return f.FileSize
}

func (f *FileInfos) Mode() fs.FileMode {
	return 0777

}

func (f *FileInfos) ModTime() time.Time {
	return f.FileModTime
}

func (f *FileInfos) IsDir() bool {
	return f.FileIsDir
}

func (f *FileInfos) Sys() any {
	return nil
}
