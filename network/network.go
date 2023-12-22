package network

import (
	"net/http"
)

type NetHandler struct {
	Mux               *http.ServeMux
	GeneralMiddleware GeneralMiddleWare
}

func NewNetwork(mux *http.ServeMux) *NetHandler {
	h := &NetHandler{Mux: mux, GeneralMiddleware: generalMiddleWare}
	return h
}
