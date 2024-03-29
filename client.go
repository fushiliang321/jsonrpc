package jsonrpc

import (
	"errors"
	"github.com/fushiliang321/jsonrpc/client"
)

type ClientInterface interface {
	Call(string, any, any, bool, any) error
	BatchAppend(string, any, any, bool, any) *error
	BatchCall() error
}

func NewClient(protocol string, ip string, port string) (ClientInterface, error) {
	switch protocol {
	case "http":
		return &client.Http{
			Ip:   ip,
			Port: port,
		}, nil
	case "tcp":
		return &client.Tcp{
			Ip:   ip,
			Port: port,
		}, nil
	}
	return nil, errors.New("The protocol can not be supported")
}
