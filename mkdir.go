package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func main() {
	path := "/Users/mac/sdk"
	infos, err := getFileInfos(path)
	if err != nil {
		return
	}
	for k, fileInfos := range infos {
		fmt.Printf("opendir path:%s", k)
		for _, v := range fileInfos {
			fmt.Println(v)
		}
	}
}

// 根据传入的文件信息得到相关的文件数据 map[string][]os.fileInfo  每个路径下的文件表示
// 期望存储的位置 map[string] 前面的map存放的是路径前缀
func getFileInfos(sourcePath string) (map[string][]FileInfos, error) {
	fileInfoMap := make(map[string][]FileInfos, 10)
	queue := NewLinkedListQueue()
	queue.Enqueue("/") //首先初始化这个前缀

	for !queue.IsEmpty() {
		dir := queue.Dequeue().(string)
		// 将这个目录下的文件夹放入
		fileInfos := fileInfoMap[dir] // 得到这个路径下的所有数据

		// 得到绝对路径
		absPath := filepath.Join(sourcePath, dir)
		fmt.Println(absPath)

		srcFd, err := os.Open(absPath)
		if err != nil {
			fmt.Println(absPath, err)
			return nil, err
		}

		for {
			readDirs, err := srcFd.Readdir(3)
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			// for each dir
			for _, readdir := range readDirs {
				relPath := filepath.Join(dir, readdir.Name()) // 得到相对路径
				if readdir.IsDir() {                          // 如果是文件夹，将对应的地址放入到队列中
					queue.Enqueue(relPath)
				} else { // 否则将文件加入到数组中
					result := FileInfoss{
						AbsPath:  filepath.Join(absPath, readdir.Name()),
						RelPath:  relPath,
						Size:     readdir.Size(),
						IsDir:    readdir.IsDir(),
						ModeTime: readdir.ModTime(),
					}
					fileInfos = append(fileInfos, result)
				}
			}
		}
		fileInfoMap[dir] = fileInfos
	}
	return fileInfoMap, nil
}

type FileInfoss struct {
	AbsPath  string // 绝对路径
	RelPath  string // 相对路径 相对前缀来说
	Size     int64
	IsDir    bool
	ModeTime time.Time
}

type node struct {
	Item interface{}
	Next *node
}

type linkedListQueue struct {
	Length int
	head   *node //头节点
	tail   *node //尾节点
}

func NewNode() *node {
	return &node{
		Item: nil,
		Next: nil,
	}
}

func NewLinkedListQueue() *linkedListQueue {
	return &linkedListQueue{
		Length: 0,
		head:   nil,
		tail:   nil,
	}
}

func (l *linkedListQueue) IsEmpty() bool {
	return l.Length == 0
}

func (l *linkedListQueue) Len() int {
	return l.Length
}

//Enqueue push
func (l *linkedListQueue) Enqueue(item interface{}) {
	buf := &node{
		Item: item,
		Next: nil,
	}
	if l.Length == 0 {
		l.tail = buf
		l.head = buf

	} else {
		l.tail.Next = buf
		l.tail = l.tail.Next
	}
	l.Length++
}

//Dequeue pop
func (l *linkedListQueue) Dequeue() (item interface{}) {
	if l.Length == 0 {
		return errors.New(
			"failed to dequeue, queue is empty")
	}

	item = l.head.Item
	l.head = l.head.Next

	// 当只有一个元素时，出列后head和tail都等于nil
	// 这时要将tail置为nil，不然tail还会指向第一个元素的位置
	// 比如唯一的元素原本为2，不做这步tail还会指向2
	if l.Length == 1 {
		l.tail = nil
	}

	l.Length--
	return
}

func (l *linkedListQueue) Traverse() (resp []interface{}) {
	buf := l.head
	for i := 0; i < l.Length; i++ {
		resp = append(resp, buf.Item, "<--")
		buf = buf.Next
	}
	return
}

func (l *linkedListQueue) GetHead() (item interface{}) {
	if l.Length == 0 {
		return errors.New(
			"failed to getHead, queue is empty")
	}
	return l.head.Item
}
