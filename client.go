package jsonrpc

import (
	"context"
	"errors"
	"github.com/fushiliang321/jsonrpc/client"
	"net"
	"net/http"
)

type ClientInterface interface {
	Call(string, any, any, bool, any) error
	CallWithContext(context.Context, string, any, any, bool, any) error
	BatchAppend(string, any, any, bool, any) *error
	BatchCall() error
	BatchCallWithContext(context.Context) error
}

func NewClient(protocol string, ip string, port string) (ClientInterface, error) {
	switch protocol {
	case "http":
		return &client.Http{
			Ip:     ip,
			Port:   port,
			Client: http.DefaultClient,
		}, nil
	case "tcp":
		return &client.Tcp{
			Ip:     ip,
			Port:   port,
			Dialer: &net.Dialer{},
		}, nil
	}
	return nil, errors.New("The protocol can not be supported")
}
