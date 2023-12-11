package main

import (
	"context"
	"fmt"
	"time"
)

// 防止GoRoutine 泄露
// 	前面的那个例子里面，goroutine还是会自己执行完，最后返回
func gen() <-chan int {
	ch := make(chan int)
	go func() {
		var n int
		for {
			ch <- n
			n++
			time.Sleep(time.Second)
		}
	}()

	return ch
}

func gen_test(ctx context.Context) <-chan int {
	ch := make(chan int)
	go func() {
		var n int
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- n:
				n++
				time.Sleep(time.Second)
			}
		}
	}()
	return ch
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	for n := range gen_test(ctx) {
		fmt.Println(n)
		if n == 5 {
			// 在break之前调用cancel函数，取消goroutine,gen函数在接受到取消信息后，直接退出
			cancelFunc()
			break
		}
	}

	// 从gen里面得到协程 得到一定的量之后直接退出？
	//for n := range gen() {
	//	fmt.Println(n)
	//	if n == 5 {
	//		break
	//	}
	//}
}
