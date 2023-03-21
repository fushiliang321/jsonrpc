package client

import (
	"fmt"
	"github.com/fushiliang321/jsonrpc/common"
	"net"
	"strconv"
	"time"
)

type Tcp struct {
	Ip          string
	Port        string
	RequestList []*common.SingleRequest
}

func (p *Tcp) BatchAppend(method string, params interface{}, result interface{}, isNotify bool, context interface{}) *error {
	singleRequest := &common.SingleRequest{
		method,
		params,
		result,
		new(error),
		isNotify,
		context,
	}
	p.RequestList = append(p.RequestList, singleRequest)
	return singleRequest.Error
}

func (p *Tcp) BatchCall() error {
	var (
		err error
		br  []interface{}
	)
	for _, v := range p.RequestList {
		var (
			req interface{}
		)
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

func (p *Tcp) Call(method string, params interface{}, result interface{}, isNotify bool, context interface{}) error {
	var (
		err error
		req []byte
	)
	if isNotify {
		req = common.JsonRs(nil, method, params, context)
	} else {
		req = common.JsonRs(strconv.FormatInt(time.Now().Unix(), 10), method, params, context)
	}
	err = p.handleFunc(req, result)
	return err
}

func (p *Tcp) handleFunc(b []byte, result interface{}) error {
	var addr = fmt.Sprintf("%s:%s", p.Ip, p.Port)
	conn, err := net.Dial("tcp", addr)
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
