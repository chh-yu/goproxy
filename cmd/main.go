package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chh-yu/goproxy/common"
	"github.com/chh-yu/goproxy/http"
	"github.com/chh-yu/goproxy/socks5"
)

func main() {
	var server common.Server
	serverType := os.Args[1]

	switch serverType {
	case "http":
		server = &http.HttpServer{
			ServerBase: common.ServerBase{
				IP:   "localhost",
				Port: 7893,
			},
		}
	case "socks5":
		server = &socks5.SOCKS5Server{
			ServerBase: common.ServerBase{
				IP:   "localhost",
				Port: 7893,
			},
			Config: &socks5.Config{
				AuthMethod: socks5.MethodNoAuth,
				TCPTimeout: 5 * time.Second,
			},
		}
	default:
		server = &socks5.SOCKS5Server{
			ServerBase: common.ServerBase{
				IP:   "localhost",
				Port: 7893,
			},
			Config: &socks5.Config{
				AuthMethod: socks5.MethodNoAuth,
				TCPTimeout: 5 * time.Second,
			},
		}
	}

	if err := server.Run(); err != nil {
		fmt.Printf("Server failed with error: %v\n", err)
	}
}
