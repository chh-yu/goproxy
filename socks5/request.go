package socks5

import (
	"fmt"
	"io"
	"net"
)

const (
	IPv4Length = 4
	IPv6length = 6
	PortLength = 2
)

type ClientRequestMessage struct {
	// Version  byte
	Cmd Command
	// Reserved byte
	AddrType AddressType
	Address  string
	Port     uint16
}
type Command = byte

const (
	CmdConnect Command = 0x01
	CmdBind    Command = 0x02
	CmdUDP     Command = 0x03
)

type AddressType = byte

const (
	TypeIPv4   AddressType = 0x01
	TypeDomain AddressType = 0x03
	TypeIPv6   AddressType = 0x04
)

type ReplyType = byte

const (
	ReplySuccess ReplyType = iota
	ReplyServerFailure
	ReplyConnectionNotAllowed
	ReplyNetworkUnreachable
	ReplyHostUnreachable
	ReplyConnectionRefused
	ReplyTTLExpired
	ReplyCommandNotSupported
	ReplyAddressTypeNotSupported
)

func NewClientRequestMessage(conn io.Reader) (*ClientRequestMessage, error) {
	// Read version, command, reserved, address type
	buf := make([]byte, 4)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	version, command, reserved, addrType := buf[0], buf[1], buf[2], buf[3]

	// Check if the fields are valid
	if version != SOCKS5Version {
		return nil, ErrVersionNotSupported
	}
	if command != CmdConnect && command != CmdBind && command != CmdUDP {
		return nil, ErrCommandNotSupported
	}
	if reserved != ReservedField {
		return nil, ErrInvalidReservedField
	}
	if addrType != TypeIPv4 && addrType != TypeDomain && addrType != TypeIPv6 {
		return nil, ErrAddressTypeNotSupported
	}

	// Read address and port
	message := ClientRequestMessage{
		Cmd:      command,
		AddrType: addrType,
	}
	switch addrType {
	case TypeIPv6:
		buf = make([]byte, IPv6length)
		fallthrough
	case TypeIPv4:
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			return nil, err
		}
		ip := net.IP(buf)
		message.Address = ip.String()
	case TypeDomain:

		if _, err := io.ReadFull(conn, buf[:1]); err != nil {
			return nil, err
		}
		domainLength := buf[0]
		if domainLength > IPv4Length {
			buf = make([]byte, domainLength)
		}
		if _, err := io.ReadFull(conn, buf[:domainLength]); err != nil {
			return nil, err
		} else {
			message.Address = string(buf[:domainLength])
		}

	}

	// Read port number
	if _, err := io.ReadFull(conn, buf[:PortLength]); err != nil {
		return nil, err
	}
	message.Port = (uint16(buf[0]) << 8) + uint16(buf[1])
	fmt.Printf("message: %v\n", message)
	return &message, nil
}

func WriteRequestSuccessMessage(conn io.Writer, ip net.IP, port uint16) error {
	// Write version, reply success, reserved, address type
	_, err := conn.Write([]byte{SOCKS5Version, ReplySuccess, ReservedField})
	if err != nil {
		return err
	}
	// Write bind IP(IPv4/IPv6)
	if _, err := conn.Write(ip); err != nil {
		return err
	}
	// Write bind port
	buf := make([]byte, 2)
	buf[0] = byte(port >> 8)
	buf[1] = byte(port)
	_, err = conn.Write(buf)

	return err
}

func WriteRequestFailureMessage(conn io.Writer, replyType ReplyType) error {
	_, err := conn.Write([]byte{SOCKS5Version, replyType, ReservedField, TypeIPv4, 0, 0, 0, 0, 0})
	return err
}
