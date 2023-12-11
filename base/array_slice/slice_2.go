package main

import "fmt"

func main() {
	s := []int{1, 2}
	s = append(s, 4, 5, 6)
	// 按照每次扩容是2倍，这个输出cap是8么
	fmt.Printf("len=%d, cap=%d", len(s), cap(s))
}
