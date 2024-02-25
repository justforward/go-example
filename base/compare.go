package main

import (
	"fmt"
	"unsafe"
)

func main() {

	//a := make([]int, 0)
	//b := make([]int, 0)
	//print(a == b) // slice can only be compared to nil

	//a := [3]int{1, 2, 3}
	//b := [3]int{1, 2, 3}
	//print(a == b) // 数组是可以直接进行比较的
	//
	//type name1 struct {
	//	n     int
	//	array unsafe.Pointer
	//}
	//
	////type name2 struct {
	////	n int
	////}
	//
	//n1 := name1{1, nil}
	//n2 := name1{2, nil}
	//print(n1 == n2)
	//
	//m := make(map[name1]int, 0)
	//print(m)

	//array1 := [3]int{1, 2, 3}
	//array2 := [3]int{1, 2, 3}
	//print(array1 == array2)
	//mm := make(map[[3]int]struct{})
	//mm[array1] = struct{}{}
	//mm[array2] = struct{}{}
	//_, ok := mm[array1]
	//print(ok)

	type Meta struct {
		Crc       uint32
		position  uint64 // 存储开始的偏移量？
		TimeStamp uint64
		KeySize   uint32
		ValueSize uint32
		Flag      uint8
	}
	fmt.Println(unsafe.Sizeof(Meta{})) // 输出的事字节数

	type Meta1 struct {
		Flag      uint8
		Crc       uint32
		KeySize   uint32
		ValueSize uint32
		position  uint64 // 存储开始的偏移量？
		TimeStamp uint64
	}

	fmt.Println(unsafe.Sizeof(Meta1{})) // 输出的事字节数

}
