package middleware

import (
	"goweb/controllers"
	"net/http"
)

// BasicAuthMiddleware ..
type BasicAuthMiddleware struct {
	Next http.Handler
}

func (bam *BasicAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if bam.Next == nil {
		bam.Next = http.DefaultServeMux
	}
	//fmt.Println(r.URL.String())
	if r.URL.String() == "/auth" {
		auth := r.Header.Get("auth")
		if auth != "1" {
			//w.WriteHeader(http.StatusUnauthorized)
			controllers.ResponseJsonError(w, r, http.StatusUnauthorized, "认证失败")
			return
		}

	}

	bam.Next.ServeHTTP(w, r)

}
