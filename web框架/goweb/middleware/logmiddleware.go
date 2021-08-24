package middleware

import (
	"log"
	"net/http"
)

// BasicAuthMiddleware ..
type LogMiddleware struct {
	Next http.Handler
}

func (lm *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if lm.Next == nil {
		lm.Next = http.DefaultServeMux
	}

	log.Printf("Access Url:%s Method:%s\n", r.URL.String(), r.Method)

	lm.Next.ServeHTTP(w, r)

}
