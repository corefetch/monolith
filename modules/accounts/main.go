package accounts

import (
	"learnt.io/core/rest"
	"learnt.io/modules/accounts/route"
)

func Service() *rest.Service {
	s := rest.NewService("accounts", "0.0.0")
	s.Post("/", route.Register)
	s.Post("/auth", route.Auth)
	s.Get("/me", route.Me)
	return s
}
