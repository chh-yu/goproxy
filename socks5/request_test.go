package socks5

import (
	"bytes"
	"net"
	"reflect"
	"testing"
)

func TestNewClientRequestMessage(t *testing.T) {
	tests := []struct {
		Version  byte
		Cmd      Command
		AddrType AddressType
		Address  []byte
		Port     []byte
		Error    error
		Message  ClientRequestMessage
	}{
		{
			Version:  SOCKS5Version,
			Cmd:      CmdConnect,
			AddrType: TypeIPv4,
			Address:  []byte{123, 35, 13, 89},
			Port:     []byte{0x00, 0x50},
			Error:    nil,
			Message: ClientRequestMessage{
				Cmd:      CmdConnect,
				AddrType: TypeIPv4,
				Address:  "123.35.13.89",
				Port:     0x0050,
			},
		},
		{
			Version:  0x00,
			Cmd:      CmdConnect,
			AddrType: TypeIPv4,
			Address:  []byte{123, 35, 13, 89},
			Port:     []byte{0x00, 0x50},
			Error:    ErrVersionNotSupported,
			Message:  ClientRequestMessage{},
		},
		{
			Version:  SOCKS5Version,
			Cmd:      CmdConnect,
			AddrType: TypeDomain,
			Address:  []byte{0x09, 0x62, 0x61, 0x69, 0x64, 0x75, 0x2e, 0x63, 0x6f, 0x6d},
			Port:     []byte{0x00, 0x50},
			Error:    nil,
			Message: ClientRequestMessage{
				Cmd:      CmdConnect,
				AddrType: TypeDomain,
				Address:  "baidu.com",
				Port:     0x0050,
			},
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		buf.Write([]byte{test.Version, test.Cmd, ReservedField, test.AddrType})
		buf.Write(test.Address)
		buf.Write(test.Port)

		message, err := NewClientRequestMessage(&buf)
		if err != test.Error {
			t.Fatalf("should get error %s, but got %s\n", test.Error, err)
		}
		if err != nil {
			return
		}
		if *message != test.Message {
			t.Fatalf("should get message %v, but got %v\n", test.Message, *message)
		}
	}

}

func TestWriteRequestSuccessMessage(t *testing.T) {
	var buf bytes.Buffer
	ip := net.IP([]byte{123, 123, 11, 11})

	err := WriteRequestSuccessMessage(&buf, ip, 1081)
	if err != nil {
		t.Fatalf("error while writing: %s", err)
	}
	want := []byte{SOCKS5Version, ReplySuccess, ReservedField, 123, 123, 11, 11, 0x04, 0x39}
	got := buf.Bytes()
	if !reflect.DeepEqual(want, buf.Bytes()) {
		t.Fatalf("message not match: want %v, got %v", want, got)
	}
}
func TestA(t *testing.T) {
	var a uint16 = 8081 // 1f71 // 0000 0001 1111 1111 : 0000 0111 0000 0001
	var b byte = byte(a - (uint16(byte(a>>8)) << 8))
	var c = byte(a)
	if b == c {
		t.Fatalf("b==c, b: %d, c:%d", b, c)
	} else {
		t.Fatalf("b!=c, b: %d, c:%d", b, c)
	}
}
