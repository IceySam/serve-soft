package network

import (
	"log"
	"net/http"
)

type GeneralMiddleWare func (handler http.Handler) http.Handler

func generalMiddleWare(handler http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("General MiddleWear")

		// do stuff

		handler.ServeHTTP(w,r)
	})
}