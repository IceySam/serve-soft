package network

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type GeneralMiddleWare func(handler http.Handler) http.HandlerFunc

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func whitelist(handler http.Handler) http.HandlerFunc  {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		res := Responses{}
		forwardedIP := r.Header.Get("X-Forwarded-For")
		clientIP := ""
		if forwardedIP != "" {
			clientIP = forwardedIP
			parts := strings.Split(forwardedIP, ",")
			clientIP = parts[0]
		} else {
			clientIP = r.RemoteAddr
		}
		pass := false
		ips := strings.Split(os.Getenv("IP_WHITELIST"), ",")
		for _, v := range ips {
			if strings.Contains(clientIP, v) {
				pass = true
				break
			}
		}

		if !pass {
			res.RespondForbidden(w, r, "IP NOT ALLOWED")
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func logging(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("[→] %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)
		handler.ServeHTTP(lrw, r)

		duration := time.Since(start)
		log.Printf("[✓] %s %s → %d (%s)", r.Method, r.RequestURI, lrw.statusCode, duration)
	})
}

func general(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := Responses{}
		secretKey := strings.TrimSpace(r.Header.Get("Secret-Key"))

		if secretKey == "" {
			res.RespondForbidden(w, r, "Secret-Key is required")
			return
		}

		if secretKey != os.Getenv("SECRET_KEY") {
			res.RespondForbidden(w, r, "Invalid Secret-Key")
			return
		}
		
		handler.ServeHTTP(w, r)
	})
}

func auth(handler http.Handler) http.HandlerFunc {
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
			return []byte(os.Getenv("JWT_SECRET")), nil
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
