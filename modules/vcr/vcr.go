package vcr

import (
	"learnt.io/core/rest"
	"learnt.io/modules/vcr/route"
)

func Service() *rest.Service {
	srv := rest.NewService("vcr", "0.0.0")
	srv.Post("/sessions", route.Create)
	srv.Get("/sessions/join", route.Join)
	return srv
}
