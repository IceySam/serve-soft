package main

import (
	"net/http"

	// "github.com/IceySam/network/network"
	// "github.com/IceySam/network/test"
)

func main() {
	mux := http.NewServeMux()

	// netHandler := network.NewNetwork(mux)

	// test.Setup(netHandler)

	http.ListenAndServe(":3000", mux)
}
