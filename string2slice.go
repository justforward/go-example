package main

import (
	"reflect"
	"unsafe"
)

/*
	要求实现slice string互转的零拷贝
   只需要共享底层 []byte 数组就可以实现 zero-copy。
*/
func main() {

}

func String2Slice(s string) []byte {

	// 字符串运行时表示
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

	sliceHeader := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}

	//
	return *(*[]byte)(unsafe.Pointer(&sliceHeader))
}
