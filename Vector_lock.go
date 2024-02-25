package main

import (
	"fmt"
	"math/rand"
	"sync"
	_ "sync/atomic"
	"time"
)

type VectorClock struct {
	id      int
	clock   []int64
	mutex   sync.Mutex
}

func NewVectorClock(id int, n int) *VectorClock {
	clock := make([]int64, n)
	clock[id] = 1
	return &VectorClock{id, clock, sync.Mutex{}}
}

func (vc *VectorClock) Update() {
	vc.mutex.Lock()
	vc.clock[vc.id] += 1
	vc.mutex.Unlock()
}

func (vc *VectorClock) Merge(other []int64) {
	vc.mutex.Lock()
	for i, v := range other {
		if v > vc.clock[i] {
			vc.clock[i] = v
		}
	}
	vc.mutex.Unlock()
}

func (vc *VectorClock) GetClock() []int64 {
	vc.mutex.Lock()
	defer vc.mutex.Unlock()
	return vc.clock
}

func (vc *VectorClock) String() string {
	return fmt.Sprintf("node%d: %v", vc.id, vc.clock)
}

func main() {
	// 创建三个节点，每个节点启动一个协程
	n := 3
	nodes := make([]*VectorClock, n)
	for i := 0; i < n; i++ {
		// 创建三个协程
		nodes[i] = NewVectorClock(i, n)
		go func(node *VectorClock) {
			for {
				// 模拟随机的事件发生
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
				node.Update()
				// 广播事件
				for _, other := range nodes {
					if other.id != node.id {
						other.Merge(node.GetClock())
					}
				}
			}
		}(nodes[i])
	}
	// 模拟运行一段时间后输出节点的时钟状态
	time.Sleep(10 * time.Second)
	for _, node := range nodes {
		fmt.Println(node)
	}
}
