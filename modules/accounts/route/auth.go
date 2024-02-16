package route

import (
	"net/http"
	"time"

	"corefetch/core"
	"corefetch/core/rest"
	"corefetch/modules/accounts/store"

	"golang.org/x/crypto/bcrypt"
)

type AuthData struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Auth(c *rest.Context) {

	var mid AuthData

	if err := c.Read(&mid); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := store.GetAccountByLogin(mid.Login)

	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(mid.Password),
	)

	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	key, err := rest.CreateKey(rest.AuthContext{
		User:   user.ID.Hex(),
		Scope:  rest.ScopeAuth,
		Expire: time.Now().Add(time.Hour),
	})

	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.Write(core.M{
		"token": key,
	})
}
