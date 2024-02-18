package rest

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

var ErrNotAuthorized = errors.New("not authorized")

type ContextKey string

var UserKey ContextKey = "user"

// Get authorization context
func GetAuthContext(c *Context) (context *AuthContext, err error) {

	var key = c.Query("access_token")

	if key == "" {
		key = c.Header("Authorization")
		if strings.IndexAny(key, "Bearer ") == 0 {
			afterBearer, _ := strings.CutPrefix(key, "Bearer ")
			key = afterBearer
		}
	}

	return ContextFromKey(key)
}

// Guard a route and authorize for specific scope
func GuardScope(scope Scope, next RestHandler) RestHandler {
	return func(c *Context) {

		ctx, err := GetAuthContext(c)

		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		if ctx.Expire.Before(time.Now()) {
			c.Status(http.StatusUnauthorized)
			return
		}

		if ctx.Scope != scope {
			c.Status(http.StatusUnauthorized)
			return
		}

		next(c)
	}
}

func GuardAuth(fn RestHandler) RestHandler {
	return GuardScope("auth", fn)
}
