package main

import "sync"

// 使用获取goroutine id
type RecursiveMutex struct {
	mu        sync.Mutex
	owner     int64
	recursion int32
}

func main() {

}
