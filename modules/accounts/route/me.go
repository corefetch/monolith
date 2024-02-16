package route

import (
	"corefetch/core/rest"
	"corefetch/modules/accounts/store"
	"net/http"
)

func Me(c *rest.Context) {

	user, err := store.GetAccount(c.User())

	if err != nil {
		c.Write(err, http.StatusNotFound)
		return
	}

	c.Write(user, http.StatusOK)
}
