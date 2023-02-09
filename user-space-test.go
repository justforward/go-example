package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hsheth2/gonet/ipv4"
	"github.com/hsheth2/gonet/tcp"
)

type contentType string

const (
	html    contentType = "text/html"
	png                 = "image/png"
	plain               = "text/plain"
	js                  = "application/javascript"
	css                 = "text/css"
	favicon             = "image/x-icon"
)

const noCache = "Cache-Control: no-cache, no-store, must-revalidate\r\n" +
	"Pragma: no-cache\r\n" +
	"Expires: 0\r\n"

var base, _ = filepath.Abs("./static")

var server *tcp.Server

func getFile(file string) ([]byte, error) {
	d, err := ioutil.ReadFile(filepath.Join(base, file))
	if err != nil {
		return nil, err
	}
	return d, nil
}

func fileType(filename string) contentType {
	fmt.Println("filetype", strings.ToLower(filepath.Ext(filename)))
	switch strings.ToLower(filepath.Ext(filename)) { // TODO more content types
	case ".html":
		return html
	case ".png":
		return png
	case ".js":
		return js
	case ".css":
		return css
	case ".ico":
		return favicon
	default:
		return plain
	}
}

func serveReq(req string) (contents []byte, tp contentType, err error) {
	d, err := getFile(req)
	if err != nil {
		// TODO more options
		n, err := getFile(filepath.Join(req, "index.html"))
		if err != nil {
			return nil, "", err
		}
		return n, html, nil
	}
	return d, fileType(req), nil
}

func respond(socket *tcp.TCB, request string) error {
	lines := strings.Split(request, "\r\n")
	fmt.Println("Request:", lines[0])

	reqLine := strings.Split(lines[0], " ")
	if strings.EqualFold(reqLine[0], "GET") {
		file := strings.Split(reqLine[1], "?")[0]
		fmt.Println(file)
		response, tp, err := serveReq(file)
		if err != nil {
			response = []byte("not found\n")
			socket.Send(
				append([]byte(
					"HTTP/1.1 404 Not Found\r\n"+
						"Content-Type: text/plain\r\n"+
						"Content-Length: "+fmt.Sprint(len(response))+"\r\n"+
						noCache+
						"Connection: close\r\n"+
						"\r\n",
				), response...),
			)
			return fmt.Errorf("serve req (finding file): %s", err)
		}
		return socket.Send(
			append([]byte(
				"HTTP/1.1 200 OK\r\n"+
					"Content-Type: "+string(tp)+"\r\n"+
					"Content-Length: "+fmt.Sprint(len(response))+"\r\n"+
					noCache+
					"Connection: close\r\n"+
					"\r\n",
			), response...),
		)
	} else {
		return errors.New(reqLine[0] + " not supported")
	}
}

func parseRespond(socket *tcp.TCB, requestFull string) (extra string, complete bool) {
	reqSplit := strings.SplitN(requestFull, "\r\n\r\n", 2)
	if len(reqSplit) <= 1 {
		return requestFull, false
	}
	request := reqSplit[0]
	extra = reqSplit[1]
	complete = true

	err := respond(socket, request)
	if err != nil {
		fmt.Println("respond", err)
	}

	return
}

// 调用链接
func connDealer(socket *tcp.TCB) {
	var buffer string
	for {

		// 定义socket接受的大小
		data, err := socket.Recv(8192) // 8kB
		if err != nil {
			if socket.IsRemoteClosed() {
				fmt.Println("Remote closed connection: closing socket")
				socket.Close()
				return
			} else {
				fmt.Println("socket recv", err)
				return
			}
		}
		buffer = buffer + string(data)
		//fmt.Println(buffer)

		//buffer = respond(socket, buffer)
		ok := true
		for ok {
			buffer, ok = parseRespond(socket, buffer)
		}
	}
}

func serverAccept() {
	for {
		// tcp accept
		socket, _, _, err := server.Accept()
		if err != nil {
			fmt.Println("tcp accept", err)
			continue
		}
		go connDealer(socket)
	}
}

func SetDir(d string) {
	base = d
}

func Run() {
	// 构建tcp的
	s, err := tcp.NewServer()
	if err != nil {
		fmt.Println("tcp server", err)
		return
	}
	// s 绑定listen的端口
	s.BindListen(80, ipv4.IPAll)
	server = s
	serverAccept()
}

func main() {
	Run()
}
