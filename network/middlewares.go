package network

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type GeneralMiddleWare func(handler http.Handler) http.Handler

func general(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("General MiddleWare")
		
		handler.ServeHTTP(w, r)
	})
}

func auth(handler http.Handler) http.Handler {
	ENV, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("General MiddleWare")
		res := Responses{}
		token := strings.TrimSpace(r.Header.Get("token"))

		if token == "" {
			res.RespondUnauthorized(w, r)
			return
		}
		c := &Claim{}
		parseClaim, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
			return []byte(ENV["JWT_SECRET"]), nil
		})

		if err != nil {
			res.RespondUnauthorized(w, r)
			return
		}

		if !parseClaim.Valid {
			res.RespondForbidden(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
