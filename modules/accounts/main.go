package accounts

import (
	"corefetch/core/rest"
	"corefetch/modules/accounts/route"
)

func Service() *rest.Service {
	srv := rest.NewService("accounts", "0.0.0")
	srv.Post("/", route.Register)
	srv.Post("/auth", route.Auth)
	srv.Get("/me", rest.GuardAuth(route.Me))
	return srv
}
