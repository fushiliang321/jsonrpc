package client

import (
	"bytes"
	"context"
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
	Client      *http.Client
	RequestList []*common.SingleRequest
}

func (p *Http) BatchAppend(method string, params any, result any, isNotify bool, contextData any) *error {
	singleRequest := &common.SingleRequest{
		Method:   method,
		Params:   params,
		Result:   result,
		Error:    new(error),
		IsNotify: isNotify,
		Context:  contextData,
	}
	p.RequestList = append(p.RequestList, singleRequest)
	return singleRequest.Error
}

func (p *Http) BatchCall(ctx context.Context) error {
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
	err = p.handleFunc(ctx, bReq, p.RequestList)
	p.RequestList = make([]*common.SingleRequest, 0)
	return err
}

func (p *Http) Call(ctx context.Context, method string, params any, result any, isNotify bool, contextData any) error {
	var (
		err error
		req []byte
	)
	if isNotify {
		req = common.JsonRs(nil, method, params, contextData)
	} else {
		req = common.JsonRs(strconv.FormatInt(time.Now().Unix(), 10), method, params, contextData)
	}
	err = p.handleFunc(ctx, req, result)
	return err
}

func (p *Http) handleFunc(ctx context.Context, b []byte, result any) error {
	var url = fmt.Sprintf("http://%s:%s", p.Ip, p.Port)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Client.Do(req)
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
