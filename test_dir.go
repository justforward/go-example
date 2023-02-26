package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {

	start := time.Now() // 获取当前时间
	path := "/share/home/linshenke/platform/job_prod/workdir/4zGg3hkCTLd.4LrN9ATvgFo/uFX_monitoringSurfaces/uFX_monitoringSurface_Monitor_Surface_1"
	//path:="/root/test_dir/dir"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	bacth := 100
	result := []FileInfo{}

	for {
		readdir, err := f.Readdir(bacth)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		for _, fi := range readdir {
			result = append(result, toFileInfos(fi))
		}
	}

	//fileInfos, err := os.ReadDir(path)
	//if err != nil {
	//	fmt.Printf("readdir error:%v", err)
	//	return
	//}
	//
	//result := []FileInfo{}
	//
	//// 得到的文件
	//for _, fileInfo := range fileInfos {
	//	fi, err := fileInfo.Info()
	//	if err != nil {
	//		return
	//	}
	//	result = append(result, toFileInfos(fi))
	//}

	elapsed := time.Since(start)
	fmt.Println("该函数执行完成耗时：", elapsed)
}

type FileInfo struct {
	Name    string
	Size    int64
	ModTime int64
	IsDir   bool
}

// FileInfo FileInfo
func toFileInfos(fi os.FileInfo) FileInfo {
	return FileInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		ModTime: fi.ModTime().UnixMilli(),
		IsDir:   fi.IsDir(),
	}
}
