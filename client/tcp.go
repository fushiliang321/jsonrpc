package client

import (
	"github.com/fushiliang321/jsonrpc/common"
	"net"
	"strconv"
	"strings"
	"time"
)

type Tcp struct {
	Ip          string
	Port        string
	RequestList []*common.SingleRequest
}

func (p *Tcp) BatchAppend(method string, params any, result any, isNotify bool, context any) *error {
	singleRequest := &common.SingleRequest{
		Method:   method,
		Params:   params,
		Result:   result,
		Error:    new(error),
		IsNotify: isNotify,
		Context:  context,
	}
	p.RequestList = append(p.RequestList, singleRequest)
	return singleRequest.Error
}

func (p *Tcp) BatchCall() error {
	var (
		err error
		br  []any
	)
	for _, v := range p.RequestList {
		var req any
		if v.IsNotify == true {
			req = common.Rs(nil, v.Method, v.Params, v.Context)
		} else {
			req = common.Rs(strconv.FormatInt(time.Now().Unix(), 10), v.Method, v.Params, v.Context)
		}
		br = append(br, req)
	}
	bReq := common.JsonBatchRs(br)
	err = p.handleFunc(bReq, p.RequestList)
	p.RequestList = make([]*common.SingleRequest, 0)
	return err
}

func (p *Tcp) Call(method string, params any, result any, isNotify bool, context any) error {
	var req []byte
	if isNotify {
		req = common.JsonRs(nil, method, params, context)
	} else {
		req = common.JsonRs(strconv.FormatInt(time.Now().Unix(), 10), method, params, context)
	}
	return p.handleFunc(req, result)
}

func (p *Tcp) handleFunc(b []byte, result any) error {
	var addrBuilder strings.Builder
	addrBuilder.WriteString(p.Ip)
	addrBuilder.WriteByte(':')
	addrBuilder.WriteString(p.Port)
	conn, err := net.Dial("tcp", addrBuilder.String())
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(b)
	if err != nil {
		return err
	}
	var buf = make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}
	err = common.GetResult(buf[:n], result)
	return err
}
