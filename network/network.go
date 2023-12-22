package network

import (
	"net/http"
)

type NetHandler struct {
	mux               *http.ServeMux
	generalMiddleware GeneralMiddleWare
}

func NewNetwork(mux *http.ServeMux) *NetHandler {
	h := &NetHandler{mux: mux, generalMiddleware: generalMiddleWare}
	return h
}
