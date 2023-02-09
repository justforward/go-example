package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {

	filePath := "/Users/mac/testSync/test1.txt"
	//write(filePath)
	insert(filePath)

}

//文件写入
func write(filePath string) {
	//os.O_CREATE:创建
	//os.O_WRONLY:只写
	//os.O_APPEND:追加
	//os.O_RDONLY:只读
	//os.O_RDWR:读写
	//os.O_TRUNC:清空

	//0644:文件的权限
	//如果没有test.txt这个文件那么就创建，并且对这个文件只进行写和追加内容。
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("文件错误,错误为:%v\n", err)
		return
	}
	defer file.Close()
	i := 0
	for i < 100 {
		runes := RandStringRunes(i)
		file.Write([]byte(runes)) //将str字符串的内容写到文件中，强制转换为byte，因为Write接收的是byte。
		i++
	}

}

func insert(filePath string) {
	i := 0
	for i < 100 {
		runes := RandStringRunes(i)
		appendStringInFile(filePath, runes)
		i++
		time.Sleep(3 * time.Second)
	}
}
func appendStringInFile(filePath, content string) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("文件打开失败: %v", err)
	}
	defer file.Close()
	// 查找文件开始插入数据
	n, _ := file.Seek(0, io.SeekStart)
	// 从末尾的偏移量开始写入内容
	_, err = file.WriteAt([]byte("\n"+content), n)
	if err != nil {
		log.Fatalf("文件写入失败: %v", err)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
