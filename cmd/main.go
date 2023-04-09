package main

import (
	"fmt"
	"log"

	"github.com/chh-yu/goproxy/socks5"
)

func main() {
	Port := 7893
	server := socks5.SOCKSServer{
		IP:   "localhost",
		Port: Port,
	}
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("run in port: %d", Port)
}
