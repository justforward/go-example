package main

import "fmt"

func main() {
	s := []int{5}
	fmt.Println(len(s)) // 1
	fmt.Println(cap(s)) // 1
	s = append(s, 7)    // 5,7
	fmt.Println(len(s)) // 2
	fmt.Println(cap(s)) // 2
	s = append(s, 9)    // 5,7,9
	fmt.Println(len(s)) // 3
	fmt.Println(cap(s)) // 4

	// 由于 s 的底层数组仍然有空间，因此并不会扩容。
	// 这样，底层数组就变成了 [5, 7, 9, 11]。注意，
	// 此时 s = [5, 7, 9]，容量为4；x = [5, 7, 9, 11]，容量为4。这里 s 不变
	x := append(s, 11) // 5,7,9,11

	// 这里还是在 s 元素的尾部追加元素，由于 s 的长度为3，容量为4，
	// 所以直接在底层数组索引为3的地方填上12。
	// 结果：s = [5, 7, 9]，y = [5, 7, 9, 12]，x = [5, 7, 9, 12]，x，y 的长度均为4，容量也均为4
	y := append(s, 12) // 5,7,9,12
	fmt.Println(s, x, y)
}
