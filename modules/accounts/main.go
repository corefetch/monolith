package accounts

import (
	"corefetch/core/rest"
	"corefetch/modules/accounts/route"
)

func Service() *rest.Service {
	s := rest.NewService("accounts", "0.0.0")
	s.Post("/", route.Register)
	s.Post("/auth", route.Auth)
	s.Get("/me", rest.GuardAuth(route.Me))
	return s
}
