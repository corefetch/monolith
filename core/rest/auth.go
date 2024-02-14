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
func GetAuthContext(i *Context) (context *AuthContext, err error) {

	var key = i.Query("access_token")

	if key == "" {
		key = i.Header("Authorization")
		if strings.IndexAny(key, "Bearer ") == 0 {
			afterBearer, _ := strings.CutPrefix(key, "Bearer ")
			key = afterBearer
		}
	}

	return ContextFromKey(key)
}

// Guard a route and authorize for specific scope
func GuardScope(scope Scope, next IntentHandler) IntentHandler {
	return func(i *Context) {

		ctx, err := GetAuthContext(i)

		if err != nil {
			i.Status(http.StatusUnauthorized)
			return
		}

		if ctx.Expire.Before(time.Now()) {
			i.Status(http.StatusUnauthorized)
			return
		}

		if ctx.Scope != scope {
			i.Status(http.StatusUnauthorized)
			return
		}

		next(i)
	}
}

func GuardAuth(fn IntentHandler) IntentHandler {
	return GuardScope("auth", fn)
}
