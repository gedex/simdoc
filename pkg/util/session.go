package util

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gedex/simdoc/pkg/datastore"
	"github.com/gedex/simdoc/pkg/model"

	"code.google.com/p/go.net/context"
	jwt "github.com/dgrijalva/jwt-go"
	webcontext "github.com/goji/context"
)

// GetUserFromRequest gets the currently authenticated user for the http.Request.
// The user details will be stored as either a simple API token or JWT bearer token.
func GetUserFromRequest(c context.Context, r *http.Request) *model.User {
	switch {
	case r.Header.Get("Authorization") != "":
		return getUserBearer(c, r)
	default:
		return nil
	}
}

// GenerateToken generates a JWT token for the user session.
func GenerateToken(c context.Context, r *http.Request, user *model.User) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["user_id"] = user.ID
	token.Claims["exp"] = time.Now().UTC().Add(time.Hour * 72).Unix()

	return token.SignedString([]byte(getSecretFromContext(c)))
}

// getUserBearer gets the currently authenticated user for the given bearer token
// (JWT).
func getUserBearer(c context.Context, r *http.Request) *model.User {
	var tokenStr = r.Header.Get("Authorization")
	fmt.Sscanf(tokenStr, "Bearer %s", &tokenStr)

	return getUserJWT(c, tokenStr)
}

// getUserJWT is a helper function that parses the user ID and retrieves the User
// data from a JWT Token.
func getUserJWT(c context.Context, token string) *model.User {
	var t, err = jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(getSecretFromContext(c)), nil
	})
	if err != nil || !t.Valid {
		return nil
	}
	var userId, ok = t.Claims["user_id"].(float64)
	if !ok {
		return nil
	}
	var user, _ = datastore.GetUserById(c, int64(userId))

	return user
}

// getSecretFromContext is a helper function to retrieve JWT secret from the
// context.
func getSecretFromContext(c context.Context) string {
	var wc = webcontext.ToC(c)

	return wc.Env["jwtSecret"].(string)
}
