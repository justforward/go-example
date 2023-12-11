package main

import (
	"fmt"
	"unsafe"
)

// unsafe point 的核心作用：
// 1）任何指针都能转化为unsafe point
//2） unsafe point都能转换成任何指针
// 3) uintptr 可以转换为unsafe point
// 4)

// Go语言中的指针不能进行计算和偏移操作，只能用来获取和修改变量的值
// go 中指针类型以及对指针的一个操作
// 1) 普通的指针类型，var m *T 指向T类型的一个普通指针
// 2) 保存指针地址的，uintptr 本质上是一个无符号类型的整数，它的大小与平台有关，本质上是保存指针地址
// 指针地址是可以进行计算的，根据这个值很容易计算出下一个指针所指向的位置
// 3) unsafe 包中提供的point，它可以指向任何类型的指针

// golang 的指针
// *类型：普通指针，只能进行读取内存存储的值，不能进行指针的运算
// unsafe.pointer 通用指针类型，用于不同指针类型之间的转换，不能进行读取内存存储的值，不能进行指针的运算
// uintptr 可以用于指针的运算,GC 不把 uintptr 当指针，uintptr 无法持有对象。uintptr 类型的目标会被回收。
// 总结：unsafe.Pointer 可以让你的变量在不同的普通指针类型转来转去，也就是表示为任意可寻址的指针类型。而 uintptr 常用于与 unsafe.Pointer 打配合，用于做指针运算。

func Counter(count *int) {
	// 支持指针类型的++ 操作 对当前指针的结果进行操作，操作的都是当前指针的内容
	*count++
}

// uintptr类型的主要是用来与unsafe.Pointer配合使用来访问和操作unsafe的内存。
//unsafe.Pointer不能执行算术操作。要想对指针进行算术运算必须这样来做：
//1.将unsafe.Pointer转换为uintptr
//2.对uintptr执行算术运算
//3.将uintptr转换回unsafe.Pointer,然后访问uintptr地址指向的对象
//需要小心的是，上面的步骤对于垃圾收集器来说应该是原子的，否则可能会导致问题。
//例如，在第1步之后，引用的对象可能被收集。如果在步骤3之后发生这种情况，指针将是一个无效的Go指针，并可能导致程序崩溃
func getSlice() {
	m := []int{0, 1, 2, 3, 4, 5}

	// 使用 uintptr 得到m中最后一个元素
	pointer := unsafe.Pointer(&m[0])
	// uintptr只能和unsafe.Pointer 之间进行类型转化
	offset := uintptr(pointer) + 5*unsafe.Sizeof(&m[0])
	value := unsafe.Pointer(offset)
	// 上面的写法可能存在风险，比如在使用过程中突然出现垃圾回收，所以必须要写成原子操作

	v := unsafe.Pointer(uintptr(unsafe.Pointer(&m[0])) + 5*unsafe.Sizeof(&m[0]))
	print(*(*int)(v))

	println(unsafe.Pointer(&m[5]))
	println(value) // 得到地址
	//println(*value)         // 直接取值，会出现问题：invalid operation: cannot indirect value (variable of type unsafe.Pointer)
	println(*(*int)(value)) // 将unsafe.pointer 转化为对应的类型才能取值

}

type Person struct {
	age  int
	name string
}

func getName() {
	p := &Person{age: 30, name: "Bob"}

	//获取到struct s中b字段的地址 先得到结构体的指针地址，然后再得到偏移
	v := unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Offsetof(p.name))

	//将其转换为一个string的指针，并且打印该指针的对应的值
	fmt.Println(*(*string)(v))
}

// 另外一个重要的要注意的是，在进行普通类型转换的时候，要注意转换的前后的类型要有相同的内存布局，
// 下面两个结构也能完成转换，就因为他们有相同的内存布局

type s1 struct {
	id   int
	name string
}

type s2 struct {
	field1 *[5]byte
	filed2 int
}

func main() {
	//count := 1
	//Counter(&count)
	//fmt.Println(count)

	getSlice()

	// 相同的内存布局也是可以的
	b := s1{name: "123"}
	var j s2
	j = *(*s2)(unsafe.Pointer(&b))
	fmt.Println(j)
}
