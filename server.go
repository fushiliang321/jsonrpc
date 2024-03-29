package jsonrpc

import (
	"errors"
	"github.com/fushiliang321/jsonrpc/server"
)

type ServerInterface interface {
	Start()
	Register(s any)
	SetBuffer(bs int)
}

func NewServer(protocol string, ip string, port string) (ServerInterface, error) {
	switch protocol {
	case "http":
		return server.NewHttpServer(ip, port), nil
	case "tcp":
		return server.NewTcpServer(ip, port), nil
	}
	return nil, errors.New("The protocol can not be supported")
}
