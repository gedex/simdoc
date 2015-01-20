package middleware

import (
	"net/http"

	"github.com/goji/context"
	"github.com/zenazn/goji/web"

	"github.com/gedex/simdoc/pkg/util"
)

// UserToContextInjector injects user information into the context.
func UserToContextInjector(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.FromC(*c)
		var user = util.GetUserFromRequest(ctx, r)
		if user != nil && user.ID != 0 {
			UserToC(c, user)
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// UserAuthorizer verifies whether current context has authenticated user information.
// If not, gives unauthorized response.
func UserAuthorizer(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if ToUser(c) == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// UserAdminAuthorizer verifies whether current context has authenticated user
// with admin role. If not, gives unauthorized or forbidden response.
func UserAdminAuthorizer(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var user = ToUser(c)
		switch {
		case user == nil:
			w.WriteHeader(http.StatusUnauthorized)
			return
		// @todo FIXME
		case user != nil && user.ID != 1:
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
