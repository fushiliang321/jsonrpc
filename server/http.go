package server

import (
	"github.com/fushiliang321/jsonrpc/common"
	"io"
	"net/http"
	"strings"
)

type Http struct {
	Ip         string
	Port       string
	Server     common.Server
	BufferSize int
}

func NewHttpServer(ip string, port string) *Http {
	return &Http{
		ip,
		port,
		common.Server{},
		BufferSize,
	}
}

func (p *Http) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", p.handleFunc)
	var addrBuilder strings.Builder
	addrBuilder.WriteString(p.Ip)
	addrBuilder.WriteByte(':')
	addrBuilder.WriteString(p.Port)
	http.ListenAndServe(addrBuilder.String(), mux)
}

func (p *Http) Register(s any) {
	p.Server.Register(s)
}

func (p *Http) SetBuffer(bs int) {
	p.BufferSize = bs
}

func (p *Http) handleFunc(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		data []byte
	)
	w.Header().Set("Content-Type", "application/json")
	if data, err = io.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	res := p.Server.Handler(data)
	w.Write(res)
}
