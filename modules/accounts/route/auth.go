package route

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"learnt.io/core"
	"learnt.io/core/rest"
	"learnt.io/modules/accounts/store"
)

type AuthData struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Auth(i *rest.Context) {

	var mid AuthData

	if err := i.Read(&mid); err != nil {
		i.Status(http.StatusBadRequest)
		return
	}

	user, err := store.GetAccountByLogin(mid.Login)

	if err != nil {
		i.Status(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(mid.Password),
	)

	if err != nil {
		i.Status(http.StatusUnauthorized)
		return
	}

	key, err := rest.CreateKey(rest.AuthContext{
		User:   user.ID.Hex(),
		Scope:  rest.ScopeAuth,
		Expire: time.Now().Add(time.Hour),
	})

	if err != nil {
		i.Status(http.StatusUnauthorized)
		return
	}

	i.Write(core.M{
		"token": key,
	})
}
