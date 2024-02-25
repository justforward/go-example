package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	// 首先创建一百个文件
	i := 1
	for i < 101 {
		filePath := "test_" + strconv.Itoa(i) + ".txt"
		fmt.Println(filePath)
		i++
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY|os.O_WRONLY, 0666)
		if err != nil {
			file.Close()
			log.Fatalf("文件打开失败: %v", err)
		}

		file.Close()
	}

	count := 1
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		time.Sleep(200 * time.Second)
		done <- true
	}()

	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			fmt.Println("Current time: ", t)
			go writeAppendFile(count)
			go truncate(count)
		}
	}

	j := 0
	for j < 10 {
		deleteFile()
		j++
	}

}

func truncate(count int) {
	filePath := "test_" + strconv.Itoa(count) + ".txt"
	f, err := os.OpenFile(filePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	s, err := f.Seek(4, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)

	// Truncate方法截取长度为size，即删除后面的内容，不管当前的偏移量在哪儿，都是从头开始截取
	// 但是其不会影响当前的偏移量
	err = f.Truncate(26)
	if err != nil {
		log.Fatal(err)
	}

}

func deleteFile() {
	intn := rand.Intn(100)
	// 往一个文件中写入数据
	filePath := "test_" + string(intn) + ".txt"
	stat, err := os.Stat(filePath)
	if err != nil {
		log.Print(stat.Name())
	}
	os.Remove(filePath)
}

func writeAppendFile(count int) {
	// 往一个文件中写入数据
	filePath := "test_" + strconv.Itoa(count) + ".txt"
	fmt.Println(filePath)
	count++
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("文件打开失败: %v", err)
	}

	defer file.Close()
	// 查找文件末尾的偏移量
	n, _ := file.Seek(0, io.SeekEnd)
	// 从末尾的偏移量开始写入内容
	i := 0
	for i < 10 {
		runes := RandStringRune(i)
		file.WriteAt([]byte(runes), n) //将str字符串的内容写到文件中，强制转换为byte，因为Write接收的是byte。
		i++
	}
	time.Sleep(time.Second)

}

var letterRune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRune(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRune[rand.Intn(len(letterRune))]
	}
	return string(b)
}
