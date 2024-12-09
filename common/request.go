package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goinggo/mapstructure"
	"reflect"
)

const (
	JsonRpc = "2.0"
)

var RequiredFields = map[string]string{
	"id":      "id",
	"jsonrpc": "jsonrpc",
	"method":  "method",
	"params":  "params",
	"context": "context",
}

type SingleRequest struct {
	Method   string
	Params   any
	Result   any
	Error    *error
	IsNotify bool
	Context  any
}

type Request struct {
	Id      string `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	Context any    `json:"context,omitempty"`
}

type NotifyRequest struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	Context any    `json:"context,omitempty"`
}

func GetRequestBody(b []byte) any {
	st := GetRequestStruct(b)
	GetStructFromJson(b, &st)
	return st
}

func GetRequestParams(b []byte, params any) error {
	Debug(reflect.TypeOf(params))
	GetStructFromJson(b, params)
	return errors.New("test")
}

func singleIndex(s *string, substr string) int {
	var (
		index = -1
		sLen  = len(*s)
	)

	for i := 0; i < sLen; i++ {
		if (*s)[i] == substr[0] {
			if index > -1 {
				return -1
			}
			index = i
		}
	}
	return index
}

func ParseRequestMethod(method string) (sName string, mName string, err error) {
	var (
		m     string
		sp    int
		first = method[0:1]
	)

	if first == "." || first == "/" {
		method = method[1:]
	}

	sp = singleIndex(&method, ".")
	if sp == -1 {
		sp = singleIndex(&method, "/")
		if sp == -1 {
			m = fmt.Sprintf("rpc: method request ill-formed: %s; need x.y or x/y", method)
			Debug(m)
			return sName, mName, errors.New(m)
		}
	}

	sName = method[:sp]
	mName = method[sp+1:]

	return sName, mName, err
}

func FilterRequestBody(jsonMap map[string]any) map[string]any {
	for k := range jsonMap {
		if _, ok := RequiredFields[k]; !ok {
			delete(jsonMap, k)
		}
	}
	return jsonMap
}

func ParseSingleRequestBody(jsonMap map[string]any) (id any, jsonrpc string, method string, params any, errCode int) {
	jsonMap = FilterRequestBody(jsonMap)
	var err error
	if _, ok := jsonMap["id"]; ok != true {
		st := NotifyRequest{}
		err = GetStruct(jsonMap, &st)
		if err != nil {
			errCode = InvalidRequest
		}
		return nil, st.JsonRpc, st.Method, st.Params, errCode
	} else {
		st := Request{}
		err = GetStruct(jsonMap, &st)
		if err != nil {
			errCode = InvalidRequest
		}
		return st.Id, st.JsonRpc, st.Method, st.Params, errCode
	}
}

func ParseRequestBody(b []byte) (any, error) {
	var (
		jsonData any
		err      = json.Unmarshal(b, &jsonData)
	)
	if err != nil {
		Debug(err)
	}
	return jsonData, err
}

func GetRequestStruct(jsonMap any) any {
	if _, ok := jsonMap.(map[string]any)["id"]; ok != true {
		return NotifyRequest{}
	} else {
		return Request{}
	}
}

func GetStructFromJson(d []byte, s any) error {
	var (
		m   string
		err error
	)
	if reflect.TypeOf(s).Kind() != reflect.Ptr {
		m = fmt.Sprintf("reflect: Elem of invalid type %s, need reflect.Ptr", reflect.TypeOf(s))
		Debug(m)
		return errors.New(m)
	}

	var jsonData any
	err = json.Unmarshal(d, &jsonData)
	if err != nil {
		Debug(err)
		return err
	}
	GetStruct(jsonData, s)
	return nil
}

func GetStruct(d any, s any) error {
	var (
		m      string
		t      reflect.Type
		typeOf = reflect.TypeOf(s)
	)
	if typeOf.Kind() != reflect.Ptr {
		m = fmt.Sprintf("reflect: Elem of invalid type %s, need reflect.Ptr", reflect.TypeOf(s))
		Debug(m)
		return errors.New(m)
	}
	t = typeOf.Elem()
	var jsonMap = make(map[string]any)
	switch d := d.(type) {
	case map[string]any:
		jsonMap = d

	case []any:
		num := t.NumField()
		if num != len(d) {
			m = fmt.Sprintf("json: The number of parameters does not match")
			Debug(m)
			return errors.New(m)
		}
		for k := 0; k < num; k++ {
			jsonMap[t.Field(k).Name] = d[k]
		}
	}

	if err := mapstructure.Decode(jsonMap, s); err != nil {
		Debug(err)
		return err
	}
	return nil
}

func Rs(id any, method string, params any, context any) any {
	var req any
	if id != nil {
		req = Request{id.(string), JsonRpc, method, params, context}
	} else {
		req = NotifyRequest{JsonRpc, method, params, context}
	}
	return req
}

func JsonRs(id any, method string, params any, context any) []byte {
	e, _ := json.Marshal(Rs(id, method, params, context))
	return e
}

func JsonBatchRs(data []any) []byte {
	e, _ := json.Marshal(data)
	return e
}
