package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 并发控制：需要接收到多个协程的信息（是否执行成功或者结束），然后进行下一步程序的执行
// 并发执行在Golang中很容易实现，只需要go func(),大多数我们会把一个大的任务拆分成多个子任务去执行，
//这时候我们就需要关心子任务是否执行成功和结束，需要收到信息进行下一步程序的执行。
//在golang 存在三种Goroutine常用的控制方式。
// 1、waitGroup
// 注意：改类型的变量也是一个值传递，当waitGroup的变量作为函数传递的时候，要作为指针才能修改对应的值

func printContext(wag *sync.WaitGroup, i int) {
	defer wag.Done()
	fmt.Println("print i:", i)

}
func waitGroup() {
	var wag sync.WaitGroup
	wag.Add(3)
	for i := 0; i < 3; i++ {
		go printContext(&wag, i)
	}
	wag.Wait()
	fmt.Println("all goroutine done")
}

//func main() {
//	waitGroup()
//}

// 2、channel
// 当一个主任务拆分为子任务去执行，子任务全部执行完毕，通过channel来通知主任务执行完毕，主任务继续向下执行。
// 比较适用于层级比较少的主任务和子任务间的通信。
// select 作为接受多个channel的

func channelAnswer() {
	ch := make(chan int)
	ok := make(chan bool)
	go son(ch, ok)
	for i := 0; i < 3; i++ {
		ch <- i
	}

	ok <- true

	fmt.Println("all is over")

}
func son(ch chan int, ok chan bool) {
	//t := time.Tick(time.Second)
	for {
		select {
		case t := <-ch:
			fmt.Println("son:", t, "run")
		case <-ok:
			fmt.Println("goroutine is over")
			break
		}

	}

}

//func main() {
//	channelAnswer()
//}

// 3、context
// 适用于嵌套多个context和组合 注意context的退出顺序？

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go foo(ctx, "Sonbar")
	fmt.Println("subwork is starting")
	time.Sleep(5 * time.Second)
	//fiveminutes is over
	cancel()
	//allgoroutine over
	time.Sleep(3 * time.Second)
	fmt.Println("main work over")
}

func foo(ctx context.Context, name string) {
	sonctx, _ := context.WithCancel(ctx)
	go boo(sonctx, name)
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "goroutine A Exit")
			return
		case <-time.After(1 * time.Second):
			fmt.Println(name, "goroutine A do something")
		}
	}

}

func boo(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "goroutine B Exit")
			return
		case <-time.After(1 * time.Second):
			fmt.Println(name, "goroutine B do something")
		}

	}
}
