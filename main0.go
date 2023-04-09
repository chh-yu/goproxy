package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	ln, err := net.Listen("tcp", ":7891")
	if err != nil {
		// 处理错误
		log.Fatalln(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			// 处理错误
			continue
		}

		// 连接到远程服务器
		remote, err := net.Dial("tcp", "47.101.139.65:80")
		if err != nil {
			// 处理错误
			conn.Close()
			continue
		}

		// 转发数据
		go func() {
			defer conn.Close()
			defer remote.Close()
			// m := make([]byte, 1024)
			// conn.Read(m)
			// fmt.Println(m)
			io.Copy(conn, remote)
		}()
		go func() {
			defer conn.Close()
			defer remote.Close()
			// m := make([]byte, 1024)
			// remote.Read(m)
			// fmt.Println(m)
			io.Copy(remote, conn)
		}()
	}

}
func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	log.Fatal(http.ListenAndServe(":7891", nil))
}
