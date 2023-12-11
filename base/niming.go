package main

// Go 语言支持匿名函数，可作为闭包
// 匿名函数是一个没有名称的函数，通常用于在函数内部定义函数，或者作为函数参数进行传递
// 特点：可以直接使用函数内部的变量，
// 匿名函数在Go语言中可以像普通变量一样被引用或者传递。
// 它们通常以函数值的形式使用，可以作为参数传递给其他函数，或者作为函数的返回值

// 1、作为函数的返回值
func getSum() func() int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}

func getSum_test(i int) func() int {

	return func() int {
		i += 1
		return i
	}

}

//func main() {
//	// 这个返回值为一个函数 getSum()
//	// 调用这个函数才生效  这个函数可以重新被调用2次
//	sum := getSum()
//	println(sum())
//	println(sum())
//
//	i := 10
//	test := getSum_test(i)
//	println(test())
//	println(test())
//}

// 2、作为函数的入参

func main() {
	add := func(i, j int) int {
		return i + j
	}

	cal := func(add func(i, j int) int, x, y int) int {
		return add(x, y)
	}

	print(cal(add, 4, 5))
}
