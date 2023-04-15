package main

import (
	"log"
	"time"

	"github.com/chh-yu/goproxy/socks5"
)

func main() {
	server := socks5.SOCKS5Server{
		IP:   "localhost",
		Port: 1080,
		Config: &socks5.Config{
			AuthMethod: socks5.MethodPassword,
			PasswordChecker: func(username, password string) bool {
				return true
			},
			TCPTimeout: 5 * time.Second,
		},
	}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
