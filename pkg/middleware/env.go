package middleware

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

func DevEnv(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if c.Env["env"] != "dev" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
