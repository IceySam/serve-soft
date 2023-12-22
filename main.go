package main

import (
	"net/http"

	"github.com/IceySam/serve-soft/network"
	"github.com/IceySam/serve-soft/test"
)

func main() {
	mux := http.NewServeMux()

	netHandler := network.NewNetwork(mux)
	
	test.Setup(netHandler)

	http.ListenAndServe(":3000", mux)
}
