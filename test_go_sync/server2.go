package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	//path := "Q3.sim"
	// 定义处理函数，读取本地的 test.html 文件，并将其作为响应返回给客户端
	http.HandleFunc("/file/content", func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("/Users/mac/Downloads/Fluent_test_Pipe.cas")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", content)
	})

	// 启动 HTTP 服务器
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("HTTP server failed: %v", err)
	}
}
