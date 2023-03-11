package jsonrpc

import (
	"errors"
	"github.com/fushiliang321/jsonrpc/client"
)

type ClientInterface interface {
	Call(string, interface{}, interface{}, bool) error
	BatchAppend(string, interface{}, interface{}, bool) *error
	BatchCall() error
}

func NewClient(protocol string, ip string, port string) (ClientInterface, error) {
	var err error
	switch protocol {
	case "http":
		return &client.Http{
			Ip:   ip,
			Port: port,
		}, err
	case "tcp":
		return &client.Tcp{
			Ip:   ip,
			Port: port,
		}, err
	}
	return nil, errors.New("The protocol can not be supported")
}
