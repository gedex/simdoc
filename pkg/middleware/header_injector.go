package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/zenazn/goji/web"
)

// OptionsRequestHeaderInjector injects HTTP header fields for OPTIONS request.
func OptionsRequestHeaderInjector(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			h := []string{"Authorization", "Content-Type"}

			if strings.HasPrefix(r.URL.Path, "/api/documents") {
				h = append(h, "Content-Range")
				h = append(h, "Content-Disposition")
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(h, ","))
			w.Header().Set("Allow", "HEAD,GET,POST,PUT,DELETE,OPTIONS")
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})

}

// HeaderInjector injects HTTP header fields into the response.
func HeaderInjector(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("X-Frame-Options", "DENY")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-XSS-Protection", "1; mode=block")
		w.Header().Add("Cache-Control", "no-cache")
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Cache-Control", "max-age=0")
		w.Header().Add("Cache-Control", "must-revalidate")
		w.Header().Add("Cache-Control", "value")
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

		h.ServeHTTP(w, r)
	})
}
