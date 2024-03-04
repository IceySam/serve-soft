package network

import (
	"net/http"
)

type NetHandler struct {
	Mux               *http.ServeMux
	Middlewares map[string]GeneralMiddleWare
}

func NewNetwork(mux *http.ServeMux, middlewares ...map[string]GeneralMiddleWare) *NetHandler {
	mid := map[string]GeneralMiddleWare { "auth": auth, "general": general}
	if len(middlewares) > 0 {
		mid = middlewares[0]
	}
	h := &NetHandler{Mux: mux, Middlewares: mid}
	return h
}
