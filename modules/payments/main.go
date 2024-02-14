package payments

import (
	"learnt.io/core"
	"learnt.io/core/rest"
)

func Service() *rest.Service {
	s := rest.NewService("payments", "0.0.0")
	s.Post("/", core.NotImplemented)
	return s
}
