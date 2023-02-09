package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var r io.Reader
	tty, err := os.OpenFile("", os.O_RDWR, 0)
	if err != nil {
		fmt.Printf("err:%v", err)
	}

	r = tty

	// 下面的断言可以执行成功
	var w io.Writer
	w = r.(io.Writer)

}
