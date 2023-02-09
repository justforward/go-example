package main

import (
	"fmt"
	"os"
)

func main() {

	filePath := ""
	for _, args := range os.Args {
		filePath = args
	}

	// 得到文件路径
	stat, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("stat file:%s,err:%v", filePath, err)
		return
	}

	linuxFileAttr := stat.Sys()
	fmt.Printf("res:%s\n, ans: %#v\n", linuxFileAttr, linuxFileAttr)

}
