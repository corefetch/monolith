package uploads

import (
	"learnt.io/core/rest"
	"learnt.io/modules/accounts/store"
)

func Service() *rest.Service {
	s := rest.NewService("uploads", "0.0.0")
	s.Post("/", rest.GuardAuth(upload))
	s.Get("/{id}", rest.GuardAuth(render))
	return s
}

func upload(c *rest.Context) {

	account, err := store.GetAccount(c.User())

	if err != nil {
		panic("account expected")
	}

	c.Write(account)
}

func render(i *rest.Context) {

}
