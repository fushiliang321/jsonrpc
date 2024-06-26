package client

import (
	"bytes"
	"fmt"
	"github.com/fushiliang321/jsonrpc/common"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Http struct {
	Ip          string
	Port        string
	RequestList []*common.SingleRequest
}

func (p *Http) BatchAppend(method string, params any, result any, isNotify bool, context any) *error {
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

func (p *Http) BatchCall() error {
	var (
		err error
		br  []any
	)
	for _, v := range p.RequestList {
		var (
			req any
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

func (p *Http) Call(method string, params any, result any, isNotify bool, context any) error {
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

func (p *Http) handleFunc(b []byte, result any) error {
	var url = fmt.Sprintf("http://%s:%s", p.Ip, p.Port)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = common.GetResult(body, result)
	return err
}
