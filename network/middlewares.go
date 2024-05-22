package network

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type GeneralMiddleWare func(handler http.Handler) http.HandlerFunc

func general(handler http.Handler) http.HandlerFunc {
	ENV, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := Responses{}
		secretKey := strings.TrimSpace(r.Header.Get("SECRET_KEY"))

		if secretKey == "" {
			res.RespondForbidden(w, r, "SECRET_KEY is required")
			return
		}

		if secretKey != ENV["SECRET_KEY"] {
			res.RespondForbidden(w, r, "Invalid SECRET_KEY")
			return
		}
		
		handler.ServeHTTP(w, r)
	})
}

func auth(handler http.Handler) http.HandlerFunc {
	ENV, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := Responses{}
		bearerToken := strings.Split(r.Header.Get("Authorization"), " ")
		if len(bearerToken) != 2 {
			res.RespondForbidden(w, r, "Malformed bearer token")
			return
		}
		if bearerToken[1] == "" {
			res.RespondUnauthorized(w, r)
			return
		}
		c := &Claim{}
		parseClaim, err := jwt.ParseWithClaims(bearerToken[1], c, func(t *jwt.Token) (interface{}, error) {
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

func ChainedMiddleware(h http.HandlerFunc, m map[string]GeneralMiddleWare) http.HandlerFunc {
	if len(m) < 1{
		return h
	}
	wrap := h
	for _, v := range m {
		wrap = v(wrap)
	}

	return wrap
}