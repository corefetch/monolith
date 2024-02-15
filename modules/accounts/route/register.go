package route

import (
	"net/http"

	"corefetch/core/rest"
	"corefetch/modules/accounts/service"
)

func Register(c *rest.Context) {

	var dataProvider service.CreateAcountData

	if err := c.Read(&dataProvider); err != nil {
		c.Write(err, http.StatusBadRequest)
		return
	}

	account, err := service.Register(dataProvider)

	if err != nil {
		c.Write(err, http.StatusInternalServerError)
		return
	}

	c.Write(account, http.StatusCreated)
}
