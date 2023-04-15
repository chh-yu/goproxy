package socks5

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var (
	ErrVersionNotSupported       = errors.New("protocol version not supported")
	ErrMethodVersionNotSupported = errors.New("sub-negotiation method version not supported")
	ErrCommandNotSupported       = errors.New("requst command not supported")
	ErrInvalidReservedField      = errors.New("invalid reserved field")
	ErrAddressTypeNotSupported   = errors.New("address type not supported")
)

const (
	SOCKS5Version = 0x05
	ReservedField = 0x00
)

type Server interface {
	Run() error
}

type SOCKS5Server struct {
	IP     string
	Port   int
	Config *Config
}

type Config struct {
	AuthMethod      Method
	PasswordChecker func(username, password string) bool
	TCPTimeout      time.Duration
}

func initConfig(config *Config) error {
	if config.AuthMethod == MethodPassword && config.PasswordChecker == nil {
		return ErrPasswordCheckerNotSet
	}
	return nil
}

func (s *SOCKS5Server) Run() error {
	// Initialize server configuration
	if err := initConfig(s.Config); err != nil {
		return err
	}

	// Listen on the specified IP:PORT
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
			if err := s.handleConnection(conn); err != nil {
				log.Printf("handle connection failure from %s: %s", conn.RemoteAddr(), err)
			}
		}()
	}
}

func (s *SOCKS5Server) handleConnection(conn net.Conn) error {
	// 协商过程
	if err := s.auth(conn); err != nil {
		return err
	}

	// Request phase
	return s.request(conn)
}

func forward(conn io.ReadWriter, targetConn io.ReadWriteCloser) error {
	defer targetConn.Close()
	go io.Copy(targetConn, conn)
	_, err := io.Copy(conn, targetConn)
	return err
}

func (s *SOCKS5Server) request(conn io.ReadWriter) error {
	// Read client request message from connection
	message, err := NewClientRequestMessage(conn)
	if err != nil {
		return err
	}

	// Check if the address type is supported
	if message.AddrType == TypeIPv6 {
		WriteRequestFailureMessage(conn, ReplyAddressTypeNotSupported)
		return ErrAddressTypeNotSupported
	}

	if message.Cmd == CmdConnect {
		return s.handleTCP(conn, message)
	} else if message.Cmd == CmdUDP {
		return s.handleUDP()
	} else {
		WriteRequestFailureMessage(conn, ReplyCommandNotSupported)
		return ErrCommandNotSupported
	}
}

func (s *SOCKS5Server) handleUDP() error {
	return nil
}

func (s *SOCKS5Server) handleTCP(conn io.ReadWriter, message *ClientRequestMessage) error {
	// 请求访问目标TCP服务
	address := fmt.Sprintf("%s:%d", message.Address, message.Port)
	targetConn, err := net.DialTimeout("tcp", address, s.Config.TCPTimeout)
	if err != nil {
		WriteRequestFailureMessage(conn, ReplyConnectionRefused)
		return err
	}

	// Send success reply
	addrValue := targetConn.LocalAddr()
	addr := addrValue.(*net.TCPAddr)
	if err := WriteRequestSuccessMessage(conn, addr.IP, uint16(addr.Port)); err != nil {
		return err
	}

	return forward(conn, targetConn)
}

func (s *SOCKS5Server) auth(conn io.ReadWriter) error {
	// Read client auth message
	clientMessage, err := NewClientAuthMessage(conn)
	if err != nil {
		return err
	}

	// Check if the auth method is supported
	var acceptable bool
	for _, method := range clientMessage.Methods {
		if method == s.Config.AuthMethod {
			acceptable = true
		}
	}
	if !acceptable {
		NewServerAuthMessage(conn, MethodNoAcceptable)
		return errors.New("method not supported")
	}
	if err := NewServerAuthMessage(conn, s.Config.AuthMethod); err != nil {
		return err
	}

	if s.Config.AuthMethod == MethodPassword {
		cpm, err := NewClientPasswordMessage(conn)
		if err != nil {
			return err
		}

		if !s.Config.PasswordChecker(cpm.Username, cpm.Password) {
			WriteServerPasswordMessage(conn, PasswordAuthFailure)
			return ErrPasswordAuthFailure
		}

		if err := WriteServerPasswordMessage(conn, PasswordAuthSuccess); err != nil {
			return err
		}
	}

	return nil
}
