package payments

import (
	"corefetch/core"
	"corefetch/core/rest"
)

func Service() *rest.Service {
	s := rest.NewService("payments", "0.0.0")
	s.Post("/", core.NotImplemented)
	return s
}
