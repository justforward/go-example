package main

import (
	"fmt"
	"time"
)

type user struct {
	name string
	age  int8
}

var u = user{name: "Ankur", age: 25}
var g = &u

func modifyUser(u *user) {
	fmt.Println("modifyUser received Value", u)
	u.name = "Anand"
}

func printUser(u <-chan *user) {
	time.Sleep(2 * time.Second)
	fmt.Println("printUser goroutine called", <-u)
}

func main() {
	c := make(chan *user, 5)
	c <- g //先把g发送到c ，根据copy value的本质，进入到chan buf里的就是一个指针地址，他还是g的值，所以打印从channel接受的元素，他就是一个&{Ankur 25} 这里并不是将指针g发送到channel中
	fmt.Println(g)
	g = &user{name: "Ankur Anand", age: 100}
	go printUser(c)
	go modifyUser(g) // 即便 g已经被修改为name: "Ankur Anand", age: 100 但是c里面接受的结果还是不变的
	time.Sleep(5 * time.Second)
	fmt.Println(g)
}
