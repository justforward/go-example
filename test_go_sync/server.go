package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	setupServer()
	time.Sleep(10 * time.Minute)
}

const path = "/Users/mac/Downloads/Fluent_test_Pipe.cas"

func handler(w http.ResponseWriter, req *http.Request) {

	// 读取文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("open %s,err:%v", path, err)
		return
	}
	fd, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("read %s, err:%v", path, err)
		return
	}
	content := bytes.NewReader(fd)
	http.ServeContent(w, req, "", time.Now(), content)
}

// set up a http server locally that will respond predictably to ranged requests
func setupServer() {
	var PORT = 8000
	s := http.NewServeMux()
	s.HandleFunc("/content", handler)
	go func() {
		var listener net.Listener
		var err error

		for {
			p := fmt.Sprintf(":%v", PORT)
			listener, err = net.Listen("tcp", p)

			if err == nil {
				break
			}
		}

		http.Serve(listener, s)
	}()
}
