package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gliderlabs/ssh"
)

func main() {

	ssh.Handle(func(s ssh.Session) {
		b, err := ioutil.ReadAll(s)
		if err != nil {
			log.Fatalf("unable to read data: %s", err)
		}
		fmt.Printf("Received data: %s", b)
	})

	err := ssh.ListenAndServe(":2233", nil,
		// 得到ssh key
		ssh.HostKeyFile("/path/to/ssh/key"),
		ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
			return password == "password"
		}),
	)
	if err != nil {
		log.Fatalf("unable to start SSH server: %s", err)
	}
}
