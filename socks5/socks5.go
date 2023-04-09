package socks5

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	ErrVersionNotSupported     = errors.New("protocol version not supported")
	ErrCommandNotSupported     = errors.New("request command not supported")
	ErrInvalidReservedField    = errors.New("invalid reserved field")
	ErrAddressTypeNotSupported = errors.New("address type not supported")
)

const (
	SOCKS5Version = 0x05
	ReservedField = 0x00
)

type Server interface {
	Run() error
}

type SOCKSServer struct {
	IP   string
	Port int
}

func (s *SOCKSServer) Run() error {
	address := fmt.Sprintf("%s:%d", s.IP, s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("connection failure from %s: %s", conn.RemoteAddr(), err)
			continue
		}
		go func() {
			defer conn.Close()
			err := handleConnection(conn)
			if err != nil {
				log.Printf("handle connection failure from %s: %s", conn.RemoteAddr(), err)
			}
		}()
	}
}
func handleConnection(conn net.Conn) error {
	// 协商过程
	if err := auth(conn); err != nil {
		return err
	}
	// 请求过程
	targetConn, err := request(conn)
	if err != nil {
		return err
	}
	// 转发过程
	return forward(conn, targetConn)
}
func forward(conn io.ReadWriter, targetConn io.ReadWriteCloser) error {
	defer targetConn.Close()
	fmt.Printf("start forward\n")
	go io.Copy(targetConn, conn)
	_, err := io.Copy(conn, targetConn)
	fmt.Printf("end forward\n")
	return err
}

func request(conn io.ReadWriter) (io.ReadWriteCloser, error) {
	message, err := NewClientRequestMessage(conn)
	if err != nil {
		return nil, err
	}
	// Check if the command is supported
	if message.Cmd != CmdConnect {
		// 返回command 不支持
		return nil, WriteRequestFailureMessage(conn, ReplyCommandNotSupported)
	}
	// Check if the address type is supported
	if message.AddrType == TypeIPv6 {
		// 返回地址类型不支持
		return nil, WriteRequestFailureMessage(conn, ReplyAddressTypeNotSupported)
	}
	// 请求访问目标TCP服务
	// message.Address:port
	address := fmt.Sprintf("%s:%d", message.Address, message.Port)
	targetConn, err := net.Dial("tcp", address)
	if err != nil {
		err := WriteRequestFailureMessage(conn, ReplyConnectionRefused)
		return nil, err
	}
	// Send success reply
	addrValue := targetConn.LocalAddr()
	addr := addrValue.(*net.TCPAddr)
	return targetConn, WriteRequestSuccessMessage(conn, addr.IP, uint16(addr.Port))
}

func auth(conn net.Conn) error {
	clientMessage, err := NewClientAuthMessage(conn)
	if err != nil {
		return err
	}
	// log.Println(clientMessage.Version, clientMessage.NMethods, clientMessage.Methods)

	// only support no-auth
	var acceptable bool
	for _, method := range clientMessage.Methods {
		if method == MethodNoAuth {
			acceptable = true
		}
	}
	if !acceptable {
		NewServerAuthMessage(conn, MethodNoAcceptable)
		return errors.New("method not supported")
	}
	return NewServerAuthMessage(conn, MethodNoAuth)
}
