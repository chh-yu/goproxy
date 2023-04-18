package main

import (
	"fmt"
	"time"

	"github.com/chh-yu/goproxy/common"
	"github.com/chh-yu/goproxy/http"
	"github.com/chh-yu/goproxy/socks5"
)

func main() {
	httpServer := &http.HttpServer{
		ServerBase: common.ServerBase{
			IP:   "localhost",
			Port: 7892,
		},
	}
	socks5Server := &socks5.SOCKS5Server{
		ServerBase: common.ServerBase{
			IP:   "localhost",
			Port: 7893,
		},
		Config: &socks5.Config{
			AuthMethod: socks5.MethodNoAuth,
			TCPTimeout: 5 * time.Second,
		},
	}

	go func() {
		if err := httpServer.Run(); err != nil {
			fmt.Printf("Server failed with error: %v\n", err)
		}
	}()

	if err := socks5Server.Run(); err != nil {
		fmt.Printf("Server failed with error: %v\n", err)
	}
}
