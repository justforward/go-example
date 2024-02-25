package custom

import (
	"fmt"
	"go-example/fdt/fs"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	bufferSize = 1024 * 1024 //1m
)

type SyncReq struct {
	start      int64
	end        int64
	chunkCount int32
	err        error
}

// SyncResp 并发写
type SyncResp struct {
	offset     int64
	length     int64
	reader     io.ReadCloser
	chunkCount int32
	err        error
}

// SyncTask 将请求的net放入到pool中
type SyncTask struct {
	targetFolder string // dest文件夹 root path
	srcFolder    string // src文件夹 root path
	srcFile      string // 源文件名，包含前置的文件夹和文件名，
	metaFile     string // 元数据 文件名
	srcFs        fs.FileSystem
	chunkSize    int32
	chunkCount   int32
	concurrency  int

	srcFileSize int64 // 文件的大小

	mx sync.RWMutex
}

func NewBoxFsSyncTask(endPoint string, concurrency int, dest string) *SyncTask {
	boxFs := fs.NewBoxFs(param.Scheme, param.Host, param.Bucket, param.ProjectID, param.Token)
	return &SyncTask{
		targetFolder: dest,
		srcFolder:    filepath.Dir(param.Path),
		srcFile:      filepath.Base(param.Path),
		srcFs:        boxFs,
		chunkSize:    12 << 20,
		concurrency:  concurrency,
	}
}
func NewSyncTask(targetFolder, srcFile, endPoint string, chunkSize int32, concurrency int) *SyncTask {
	networkFs := fs.NewNetworkFs(endPoint)
	return &SyncTask{
		targetFolder: targetFolder,
		srcFile:      srcFile,
		srcFs:        networkFs,
		chunkSize:    chunkSize,
		concurrency:  concurrency,
	}
}

func (req *SyncTask) Sync() {
	log.Printf("blocksize:%d", req.chunkSize)
	now := time.Now()
	syncReq := make(chan *SyncReq)
	metaChunk := int32(0)
	go req.productReq(syncReq, &metaChunk)

	syncResp := make(chan *SyncResp)
	go req.productResp(syncReq, syncResp)

	req.writeToFile(syncResp, &metaChunk)

	fmt.Println(fmt.Sprintf("metaChunk total count:%d", metaChunk))
	if req.chunkCount == metaChunk {
		metaFile := filepath.Join(req.targetFolder, req.srcFile+"_.meta._")
		if err := os.Remove(metaFile); err != nil {
			panic(err)
		}
	}
	sub := time.Now().Sub(now)
	fmt.Println(fmt.Sprintf("sync use time:%s, avg speed:%f MB/s", sub, float64(req.srcFileSize/1024/1024)/sub.Seconds()))

}

func (req *SyncTask) writeToFile(resp <-chan *SyncResp, metaChunk *int32) {

	metaFileName := filepath.Join(req.targetFolder, req.srcFile+"_.meta._")
	metaFile, err := os.OpenFile(metaFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(metaFile)

	// 写入的文件内容
	wg := sync.WaitGroup{}

	for i := 1; i <= req.concurrency; i++ {
		wg.Add(1)
		go req.respConsumer(resp, metaFile, metaChunk, &wg)
	}

	wg.Wait()
}

func (req *SyncTask) respConsumer(resp <-chan *SyncResp, metaFile *os.File, metaChunk *int32, wg *sync.WaitGroup) {
	defer wg.Done()
	for info := range resp {
		targetFile, err := os.OpenFile(filepath.Join(req.targetFolder, req.srcFile), os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}

		if _, err := targetFile.Seek(info.offset, io.SeekStart); err != nil {
			panic(err)
		}

		if n, err := io.Copy(targetFile, io.LimitReader(info.reader, info.length)); err != nil {
			log.Printf("respConsumer write:%d,offset:%d,length:%d", n, info.offset, info.length)
			panic(err)
		}

		info.reader.Close()

		// 每次记录进度和速度

		if _, err := metaFile.Seek(int64(info.chunkCount), io.SeekStart); err != nil {
			panic(err)
		}

		if _, err := metaFile.Write([]byte(string(rune(1)))); err != nil {
			panic(err)
		}

		//req.mx.Lock()
		//*metaChunk++
		//req.mx.Unlock()

		targetFile.Close()
	}
}

func (req *SyncTask) productReq(syncReq chan<- *SyncReq, metaChuck *int32) {
	defer close(syncReq)

	// 请求得到文件的信息
	srcFile, err := req.srcFs.Stat(filepath.Join(req.srcFolder, req.srcFile))
	if err != nil {
		panic(err)
	}

	srcSize := srcFile.Size()
	count := int32(srcSize / int64(req.chunkSize))
	if (srcSize % int64(req.chunkSize)) != 0 {
		count++
	}
	req.chunkCount = count
	req.srcFileSize = srcSize

	req.MetaFileNotExist(syncReq, count, srcSize)

	//metaFilePath := filepath.Join(req.targetFolder, req.srcFile+metaFileSuffix)
	//destFilePath := filepath.Join(req.targetFolder, req.srcFile)

	//if err = os.MkdirAll(filepath.Dir(metaFilePath), 0755); err != nil {
	//	panic(err)
	//}

	//  打开dest的文件
	//_, err = os.Stat(destFilePath)
	//if err != nil && os.IsNotExist(err) {
	//	metaFile, err := os.Create(metaFilePath)
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer func(file *os.File) {
	//		err := file.Close()
	//		if err != nil {
	//			panic(err)
	//		}
	//	}(metaFile)
	//
	//	 req.MetaFileNotExist(syncReq, count, srcSize)
	//}
	//
	//if err == nil {
	//	metaFd, err := os.Open(metaFilePath)
	//	if err != nil && !os.IsNotExist(err) {
	//		syncReq <- &SyncReq{err: err}
	//		return
	//	}
	//	if err != nil && os.IsNotExist(err) {
	//		return
	//	}
	//
	//	content, err := io.ReadAll(metaFd)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	number, err := strconv.Atoi(string(content))
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	req.MetaFileExist(syncReq, count, int32(number), srcSize, metaChuck)
	//}

}

// MetaFileNotExist 第一构建的时候
func (req *SyncTask) MetaFileNotExist(syncReq chan<- *SyncReq, count int32, fileSize int64) {
	index := int32(0)
	for index < count {
		syncReq <- &SyncReq{
			start:      int64(index) * int64(req.chunkSize),
			end:        int64(math.Min(float64(int64(index+1)*int64(req.chunkSize))-1, float64(fileSize-1))),
			chunkCount: index,
		}
		index++
	}
}
func (req *SyncTask) MetaFileExist(syncReq chan<- *SyncReq, count, number int32, fileSize int64, metaChuck *int32) {
	index := int32(0)
	for index < count { // 从0开始的下标
		// 如果当前为不为1 需要进行传递的参数，需要进行写的数据
		if !isBitSet(uint64(number), uint64(index)) {
			syncReq <- &SyncReq{
				start:      int64(index) * int64(req.chunkSize),
				end:        int64(math.Min(float64(int64(index+1)*int64(req.chunkSize)), float64(fileSize))),
				chunkCount: index,
			}
		} else {
			*metaChuck++
		}
		index++
	}
}

func (req *SyncTask) productResp(syncReq <-chan *SyncReq, syncResp chan<- *SyncResp) {
	defer close(syncResp)

	wg := sync.WaitGroup{}
	for i := 1; i <= req.concurrency; i++ {
		wg.Add(1)
		go req.reqConsumer(syncReq, syncResp, &wg)
	}
	wg.Wait()
}

func (req *SyncTask) reqConsumer(syncReq <-chan *SyncReq, syncResp chan<- *SyncResp, wg *sync.WaitGroup) {
	defer wg.Done()
	for infos := range syncReq {
		tries := 0
		for tries < 3 {
			if output, err := req.RemoteReq(infos, req.srcFs, filepath.Join(req.srcFolder, req.srcFile)); err == nil {
				syncResp <- output
				break
			}
			tries++
		}
	}
}

func (req *SyncTask) RemoteReq(input *SyncReq, fs fs.FileSystem, path string) (output *SyncResp, err error) {
	resp, err := fs.Download(path, input.start, input.end)
	if err != nil {
		panic(err)
	}

	return &SyncResp{
		offset:     input.start,
		length:     input.end - input.start + 1,
		reader:     resp,
		chunkCount: input.chunkCount,
	}, nil
}

// bit
func isBitSet(num, position uint64) bool {
	shifted := num >> position
	return (shifted & 1) == 1
}
