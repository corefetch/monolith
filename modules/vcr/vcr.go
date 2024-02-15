package vcr

import (
	"corefetch/core/rest"
	"corefetch/modules/vcr/route"
)

func Service() *rest.Service {
	srv := rest.NewService("vcr", "0.0.0")
	srv.Post("/sessions", route.Create)
	srv.Get("/sessions/join", route.Join)
	return srv
}
